import { useEffect, useRef, useState, useCallback } from 'react';
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';
import { useMutation } from '@tanstack/react-query';
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from '@/components/ui/form.tsx';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select.tsx';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { useAttendanceCategories } from "@/query/attendance-category.query.ts";
import { useOffice } from "@/query/useOffice.ts";
import Loading from "@/components/Loading.tsx";
import { Input } from '@/components/ui/input';
import { createAttendanceRecord } from "@/apis/attendance-record.api.ts";
import { toast } from "sonner";
import axios from "axios";
import { toVietnamISOString } from "@/lib/utils.ts";

// Fix Leaflet icon issue
delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
    iconRetinaUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-icon-2x.png',
    iconUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-icon.png',
    shadowUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-shadow.png',
});

const formSchema = z.object({
    category_id: z.string({
        required_error: 'Vui lòng chọn hình thức làm việc',
    }),
    office_id: z.string().min(1, 'Vui lòng chọn văn phòng'),
    note_request: z.string().optional(),
    image: z
        .instanceof(File, {
            message: 'Vui lòng chụp ảnh chấm công',
        }),
    gps_location: z.string(),
    timestamp: z.string().datetime()
});

const AttendancePage = () => {
    const { data: categories, isLoading: isCategoriesLoading } = useAttendanceCategories()
    const [position, setPosition] = useState<[number, number] | null>(null);
    const [cameraError, setCameraError] = useState<string | null>(null);

    const [currentTime, setCurrentTime] = useState(new Date());
    const videoRef = useRef<HTMLVideoElement>(null);
    const canvasRef = useRef<HTMLCanvasElement>(null);
    const streamRef = useRef<MediaStream | null>(null);

    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            category_id: '',
            office_id: '',
            note_request: '',
            image: undefined,
            gps_location: '',
            timestamp: new Date().toISOString()
        },
    });

    const { data: offices = [], isLoading: isOfficesLoading } = useOffice();

    const [cameraReady, setCameraReady] = useState(false);

    // Initialize camera
    const initCamera = useCallback(async () => {
        try {
            // Check if camera is already initialized
            if (streamRef.current) {
                return;
            }

            // Check browser support
            if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
                setCameraError("Trình duyệt không hỗ trợ camera");
                toast.error("Trình duyệt không hỗ trợ camera");
                return;
            }

            // Request camera permissions
            streamRef.current = await navigator.mediaDevices.getUserMedia({
                video: {
                    facingMode: 'user',
                    width: { ideal: 640 },
                    height: { ideal: 480 }
                }
            });

            if (videoRef.current) {
                videoRef.current.srcObject = streamRef.current;
                videoRef.current.onloadedmetadata = () => {
                    setCameraReady(true);
                    videoRef.current?.play().catch(err => {
                        console.error("Video play error:", err);
                        setCameraError("Không thể phát video từ camera");
                    });
                };
            }
        } catch (err) {
            console.error("Camera error:", err);
            const errorMessage = err instanceof Error ? err.message : 'Lỗi không xác định';
            setCameraError(`Không thể mở camera: ${errorMessage}`);
            toast.error(`Không thể mở camera: ${errorMessage}`);
        }
    }, []);

    useEffect(() => {
        initCamera();

        // Cleanup function
        return () => {
            if (streamRef.current) {
                streamRef.current.getTracks().forEach(track => track.stop());
                streamRef.current = null;
            }
        };
    }, [initCamera]);

    const captureImage = useCallback(async (): Promise<File | null> => {
        if (!canvasRef.current || !videoRef.current || !cameraReady) {
            toast.error("Camera chưa sẵn sàng");
            return null;
        }

        try {
            const context = canvasRef.current.getContext("2d");
            if (!context) {
                toast.error("Không thể khởi tạo canvas");
                return null;
            }

            // Set canvas dimensions to match video
            canvasRef.current.width = videoRef.current.videoWidth;
            canvasRef.current.height = videoRef.current.videoHeight;

            // Draw current video frame to canvas
            context.drawImage(
                videoRef.current,
                0, 0,
                canvasRef.current.width,
                canvasRef.current.height
            );

            // Convert to blob and create file
            return new Promise((resolve) => {
                canvasRef.current?.toBlob((blob) => {
                    if (blob) {
                        const file = new File([blob], `attendance_${Date.now()}.jpg`, {
                            type: "image/jpeg"
                        });
                        resolve(file);
                    } else {
                        resolve(null);
                    }
                }, "image/jpeg", 0.8); // 80% quality
            });
        } catch (error) {
            console.error("Capture error:", error);
            toast.error("Lỗi khi chụp ảnh");
            return null;
        }
    }, [cameraReady]);

    const submitAttendanceMutation = useMutation({
        mutationFn: async (data: z.infer<typeof formSchema>) => {
            const formData = new FormData();
            formData.append('category_id', data.category_id);
            formData.append('office_id', data.office_id);
            formData.append('note_request', data.note_request || '');
            formData.append('gps_location', data.gps_location);
            formData.append('timestamp', data.timestamp);
            formData.append('image', data.image);
            return await createAttendanceRecord(formData);
        },
        onSuccess: () => {
            toast.success('Chấm công thành công!');
            form.reset();
        },
        onError: (error) => {
            if (axios.isAxiosError(error)) {
                toast.error(error?.response?.data.message || 'Lỗi khi chấm công');
            } else {
                toast.error('Lỗi khi chấm công');
            }
        }
    });

    const onSubmit = async (data: z.infer<typeof formSchema>) => {
        try {
            const submissionData = {
                ...data,
                timestamp: toVietnamISOString(),
                gps_location: position ? `${position[0]},${position[1]}` : ''
            };

            await submitAttendanceMutation.mutateAsync(submissionData);
        } catch (error) {
            console.error('Submission error:', error);
        }
    };

    const handleSubmit = async () => {
        const imageFile = await captureImage();
        if (!imageFile) {
            return;
        }

        form.setValue("image", imageFile);
        form.handleSubmit(onSubmit)();
    };

    const calculateDistance = useCallback((
        lat1: number,
        lon1: number,
        lat2: number,
        lon2: number
    ) => {
        const R = 6371; // km
        const dLat = ((lat2 - lat1) * Math.PI) / 180;
        const dLon = ((lon2 - lon1) * Math.PI) / 180;
        const a =
            Math.sin(dLat / 2) ** 2 +
            Math.cos((lat1 * Math.PI) / 180) *
            Math.cos((lat2 * Math.PI) / 180) *
            Math.sin(dLon / 2) ** 2;
        const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
        return R * c;
    }, []);

    // Geolocation effect
    useEffect(() => {
        if ('geolocation' in navigator) {
            const geoOptions = {
                enableHighAccuracy: true,
                timeout: 10000,
                maximumAge: 60000
            };

            navigator.geolocation.getCurrentPosition(
                (pos) => {
                    const newPosition: [number, number] = [
                        pos.coords.latitude,
                        pos.coords.longitude
                    ];
                    setPosition(newPosition);
                    form.setValue('gps_location', `${newPosition[0]},${newPosition[1]}`);
                },
                (err) => {
                    console.error('Geolocation error:', err);
                    toast.error('Không thể lấy vị trí. Vui lòng kiểm tra quyền truy cập vị trí.');
                },
                geoOptions
            );
        } else {
            toast.error('Trình duyệt không hỗ trợ geolocation');
        }
    }, [form]);

    // Timer effect
    useEffect(() => {
        const timer = setInterval(() => {
            const now = new Date();
            setCurrentTime(now);
            form.setValue('timestamp', now.toISOString());
        }, 1000);

        return () => clearInterval(timer);
    }, [form]);

    // Auto-select nearest office effect
    useEffect(() => {
        if (position && offices.length > 0) {
            const nearestOffice = offices.reduce((nearest, office) => {
                const d1 = calculateDistance(
                    position[0],
                    position[1],
                    office.latitude,
                    office.longitude
                );
                const d2 = calculateDistance(
                    position[0],
                    position[1],
                    nearest.latitude,
                    nearest.longitude
                );
                return d1 < d2 ? office : nearest;
            }, offices[0]);

            if (nearestOffice) {
                form.setValue("office_id", nearestOffice.office_id);
            }
        }
    }, [position, offices, form, calculateDistance]);

    const formattedTime = (() => {
        const d = currentTime;
        const weekdays = [
            'Chủ Nhật',
            'Thứ 2',
            'Thứ 3',
            'Thứ 4',
            'Thứ 5',
            'Thứ 6',
            'Thứ 7',
        ];
        return `Thời gian: ${d.toLocaleTimeString('vi-VN', {
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit'
        })} - ${weekdays[d.getDay()]} ngày ${d
            .getDate()
            .toString()
            .padStart(2, '0')}/${(d.getMonth() + 1)
            .toString()
            .padStart(2, '0')}/${d.getFullYear()}`;
    })();

    const sortedOffices = position
        ? [...offices].sort((a, b) => {
            const da = calculateDistance(
                position[0],
                position[1],
                a.latitude,
                a.longitude
            );
            const db = calculateDistance(
                position[0],
                position[1],
                b.latitude,
                b.longitude
            );
            return da - db;
        })
        : offices;

    if (isCategoriesLoading || isOfficesLoading) {
        return <Loading />;
    }

    return (
        <div className="container mx-auto p-4">
            <h1 className="text-3xl font-bold mb-6">Chấm công</h1>
            <div className="p-6 space-y-6 bg-white rounded-2xl shadow-lg">
                {/* Bản đồ */}
                <div className="h-[300px] w-full rounded-lg overflow-hidden">
                    {position ? (
                        <MapContainer
                            center={position}
                            zoom={16}
                            style={{ height: '100%', width: '100%' }}
                            zoomControl={true}
                        >
                            <TileLayer
                                attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                            />
                            <Marker
                                position={position}
                                icon={L.icon({
                                    iconUrl: "https://unpkg.com/leaflet@1.7.1/dist/images/marker-icon.png",
                                    iconSize: [25, 41],
                                    iconAnchor: [12, 41]
                                })}
                            >
                                <Popup>Vị trí hiện tại của bạn</Popup>
                            </Marker>
                        </MapContainer>
                    ) : (
                        <div className="flex items-center justify-center h-full">
                            <p>Đang lấy vị trí...</p>
                        </div>
                    )}
                </div>

                {/* Camera section */}
                <div className="flex flex-col items-center">
                    <video
                        ref={videoRef}
                        autoPlay
                        muted
                        playsInline
                        className="rounded-lg border border-gray-300 mx-auto"
                        style={{ width: 300, height: 300, background: '#000' }}
                    />

                    {!cameraReady && !cameraError && (
                        <p className="mt-2 text-gray-500">Đang khởi động camera...</p>
                    )}

                    {cameraError && (
                        <p className="mt-2 text-red-500">{cameraError}</p>
                    )}
                </div>

                {/* Hidden canvas for image capture */}
                <canvas ref={canvasRef} className="hidden" />

                <div className="text-center text-lg font-medium">
                    {formattedTime}
                </div>

                <Form {...form}>
                    <form className="space-y-4">
                        <FormField
                            control={form.control}
                            name="category_id"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Hình thức làm việc</FormLabel>
                                    <Select
                                        onValueChange={field.onChange}
                                        defaultValue={field.value}
                                    >
                                        <FormControl>
                                            <SelectTrigger className="w-full">
                                                <SelectValue placeholder="Chọn hình thức làm việc" />
                                            </SelectTrigger>
                                        </FormControl>
                                        <SelectContent>
                                            {categories?.map((item) => (
                                                <SelectItem
                                                    key={item.attendance_category_id}
                                                    value={item.attendance_category_id}
                                                >
                                                    {item.attendance_category_name}
                                                </SelectItem>
                                            ))}
                                        </SelectContent>
                                    </Select>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />

                        <FormField
                            control={form.control}
                            name="office_id"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Văn phòng</FormLabel>
                                    <Select
                                        onValueChange={field.onChange}
                                        value={field.value}
                                    >
                                        <FormControl>
                                            <SelectTrigger className="w-full">
                                                <SelectValue placeholder="Chọn văn phòng" />
                                            </SelectTrigger>
                                        </FormControl>
                                        <SelectContent>
                                            {sortedOffices.map((office) => (
                                                <SelectItem
                                                    key={office.office_id}
                                                    value={office.office_id}
                                                >
                                                    {office.office_name}
                                                </SelectItem>
                                            ))}
                                        </SelectContent>
                                    </Select>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />

                        <FormField
                            control={form.control}
                            name="note_request"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Yêu cầu/Ghi chú</FormLabel>
                                    <FormControl>
                                        <Input
                                            placeholder="Nhập yêu cầu hoặc ghi chú"
                                            {...field}
                                        />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />

                        {/* Hidden fields */}
                        <FormField
                            control={form.control}
                            name="gps_location"
                            render={({ field }) => (
                                <FormItem className="hidden">
                                    <FormControl>
                                        <Input type="hidden" {...field} />
                                    </FormControl>
                                </FormItem>
                            )}
                        />

                        <FormField
                            control={form.control}
                            name="timestamp"
                            render={({ field }) => (
                                <FormItem className="hidden">
                                    <FormControl>
                                        <Input type="hidden" {...field} />
                                    </FormControl>
                                </FormItem>
                            )}
                        />

                        <FormField
                            control={form.control}
                            name="image"
                            render={({ field }) => (
                                <FormItem className="hidden">
                                    <FormControl>
                                        <Input type="hidden" {...field} />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />

                        <button
                            type="button"
                            onClick={handleSubmit}
                            className="w-full mt-4 px-4 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed"
                            disabled={submitAttendanceMutation.isPending || !cameraReady || !!cameraError}
                        >
                            {submitAttendanceMutation.isPending ? (
                                <span className="flex items-center justify-center">
                  <Loading />
                  <span className="ml-2">Đang xử lý...</span>
                </span>
                            ) : (
                                'Chấm công'
                            )}
                        </button>
                    </form>
                </Form>
            </div>
        </div>
    );
};

export default AttendancePage;