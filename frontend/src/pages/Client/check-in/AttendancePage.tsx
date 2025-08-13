import {useEffect, useRef, useState} from 'react';
import {MapContainer, TileLayer, Marker, Popup} from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';
import {useMutation} from '@tanstack/react-query';
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
import {useForm} from 'react-hook-form';
import {z} from 'zod';
import {zodResolver} from '@hookform/resolvers/zod';
import {useAttendanceCategories} from "@/query/attendance-category.query.ts";
import {useOffice} from "@/query/useOffice.ts";
import Loading from "@/components/Loading.tsx";
import { Input } from '@/components/ui/input';
import {createAttendanceRecord} from "@/apis/attendance-record.api.ts";
import {toast} from "sonner";
import axios from "axios";

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
    const {data: categories, isLoading: isCategoriesLoading} = useAttendanceCategories()
    const [cameraReady, setCameraReady] = useState(false);
    const [position, setPosition] = useState<[number, number] | null>(null);
    const [cameraError, setCameraError] = useState('');
    const [currentTime, setCurrentTime] = useState(new Date());
    const [imageData, setImageData] = useState<string | null>(null);
    const videoRef = useRef<HTMLVideoElement>(null);
    const canvasRef = useRef<HTMLCanvasElement>(null);
    const [isCamera, setIsCamera] = useState(true); // Thêm state này để kiểm soát việc sử dụng camera
    const streamRef = useRef<MediaStream | null>(null); // Sử dụng useRef để lưu stream
    const [cameraPermission, setCameraPermission] = useState<'prompt' | 'granted' | 'denied'>('prompt');
    const [isCameraAvailable, setIsCameraAvailable] = useState(true);

    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            category_id: '',
            office_id: '',
            note_request: '',
            image: '',
            gps_location: '',
            timestamp: new Date().toISOString()
        },
    });

    const {data: offices = [], isLoading: isOfficesLoading} = useOffice();

    const submitAttendanceMutation = useMutation({
        mutationFn: async (data: z.infer<typeof formSchema>) => {
            const formData = new FormData();
            formData.append('category_id', data.category_id);
            formData.append('office_id', data.office_id);
            formData.append('note_request', data.note_request || '');
            formData.append('gps_location', data.gps_location);
            formData.append('timestamp', data.timestamp);
            formData.append('image', data.image); // File here!

            return await createAttendanceRecord(formData);
        },
        onSuccess: () => {
            toast.success('Chấm công thành công!');
            form.reset();
            setImageData(null);
        },
        onError: (error) => {
            if (axios.isAxiosError(error)) {
                console.log(error.response);
                toast.error(error.response.data.message);
            }
        }
    });


    const onSubmit = async (data: z.infer<typeof formSchema>) => {
        try {
            const submissionData = {
                ...data,
                timestamp: new Date().toISOString(),
                gps_location: position ? `${position[0]},${position[1]}` : ''
            };

            await submitAttendanceMutation.mutateAsync(submissionData);
        } catch (error) {
            console.error('Submission error:', error);
        }
    };

    const calculateDistance = (
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
    };

    useEffect(() => {
        if ('geolocation' in navigator) {
            navigator.geolocation.getCurrentPosition(
                (pos) => {
                    setPosition([pos.coords.latitude, pos.coords.longitude]);
                    form.setValue('gps_location', `${pos.coords.latitude},${pos.coords.longitude}`);
                },
                (err) => {
                    console.error('Geolocation error:', err);
                }
            );
        }
    }, [form]);

    useEffect(() => {
        const timer = setInterval(() => {
            const now = new Date();
            setCurrentTime(now);
            form.setValue('timestamp', now.toISOString());
        }, 1000);
        return () => clearInterval(timer);
    }, [form]);

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
        return `Thời gian: ${d.toLocaleTimeString([], {
            hour: '2-digit',
            minute: '2-digit',
        })}     ${weekdays[d.getDay()]} ngày ${d
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

    useEffect(() => {
        if (position && offices.length > 0) {
            const nearestOffice = offices.reduce((nearest, office) => {
                const d1 = calculateDistance(position[0], position[1], office.latitude, office.longitude);
                const d2 = calculateDistance(position[0], position[1], nearest.latitude, nearest.longitude);
                return d1 < d2 ? office : nearest;
            }, offices[0]);

            if (nearestOffice) {
                console.log("Auto-select office:", nearestOffice.office_name);
                form.setValue("office_id", nearestOffice.office_id);
            }
        }
    }, [position, offices, form]);

    useEffect(() => {
        const initCamera = async () => {
            try {
                // Kiểm tra quyền truy cập camera trước
                const permissionStatus = await navigator.permissions.query({ name: 'camera' as any });
                setCameraPermission(permissionStatus.state);

                if (permissionStatus.state === 'denied') {
                    setCameraError('Quyền truy cập camera bị từ chối. Vui lòng cấp quyền trong cài đặt trình duyệt.');
                    setIsCameraAvailable(false);
                    return;
                }

                // Kiểm tra thiết bị có camera không
                const devices = await navigator.mediaDevices.enumerateDevices();
                const hasCamera = devices.some(device => device.kind === 'videoinput');

                if (!hasCamera) {
                    setCameraError('Không tìm thấy camera trên thiết bị');
                    setIsCameraAvailable(false);
                    return;
                }

                // Khởi tạo camera
                const mediaStream = await navigator.mediaDevices.getUserMedia({
                    video: {
                        width: { ideal: 1280 },
                        height: { ideal: 720 },
                        facingMode: 'user'
                    }
                });

                streamRef.current = mediaStream;

                if (videoRef.current) {
                    videoRef.current.srcObject = mediaStream;
                    videoRef.current.onloadedmetadata = () => {
                        videoRef.current?.play();
                        setCameraReady(true);
                    };
                }
            } catch (err) {
                console.error("Lỗi camera:", err);
                setCameraError('Không thể khởi động camera. Vui lòng kiểm tra quyền truy cập hoặc thử lại.');
                setIsCameraAvailable(false);
            }
        };

        if (isCamera && isCameraAvailable) {
            initCamera();
        }

        return () => {
            if (streamRef.current) {
                streamRef.current.getTracks().forEach(track => {
                    track.stop();
                });
                streamRef.current = null;
            }
        };
    }, [isCamera, isCameraAvailable]);

    if (isCategoriesLoading || isOfficesLoading) {
        return <Loading />;
    }

    const handleSubmit = async () => {
        if (videoRef.current && canvasRef.current) {
            const ctx = canvasRef.current.getContext('2d');
            if (!ctx) {
                toast.error('Không thể xử lý ảnh từ camera.');
                return;
            }

            ctx.drawImage(videoRef.current, 0, 0, 320, 240);
            canvasRef.current.toBlob((blob) => {
                if (!blob) {
                    toast.error('Không thể chuyển ảnh sang File.');
                    return;
                }

                const file = new File([blob], `attendance-${Date.now()}.png`, {
                    type: 'image/png',
                });

                setImageData(URL.createObjectURL(file)); // preview nếu cần
                form.setValue('image', file, { shouldValidate: true });

                // Gọi submit sau khi đã có file
                form.handleSubmit(onSubmit)();
            }, 'image/png');
        } else {
            toast.error('Camera không sẵn sàng.');
        }
    };


    return (
        <div>
            <h1 className="text-3xl font-bold">Chấm công</h1>
            <div className="p-4 space-y-6 h-full flex flex-col items-center justify-center bg-white rounded-2xl">
                {/* Bản đồ */}
                <div className="h-[300px] w-[350px]">
                    {position ? (
                        <MapContainer
                            center={position}
                            zoom={16}
                            style={{height: '100%', width: '100%'}}
                        >
                            <TileLayer
                                attribution='&copy; <a href="http://osm.org">OpenStreetMap</a> contributors'
                                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                            />
                            <Marker position={position} icon={L.icon({
                                iconUrl: "https://unpkg.com/leaflet@1.7.1/dist/images/marker-icon.png",
                                iconSize: [25, 41],
                                iconAnchor: [12, 41]
                            })}>
                                <Popup>Vị trí hiện tại của bạn</Popup>
                            </Marker>
                        </MapContainer>
                    ) : (
                        <p>Đang lấy vị trí...</p>
                    )}
                </div>

                {/* Camera */}
                <div>
                    {isCamera && isCameraAvailable ? (
                        <div>
                            {!cameraReady && !cameraError && (
                                <div className="w-[320px] h-[240px] flex items-center justify-center bg-black rounded-lg">
                                    <p className="text-white">Đang khởi tạo camera...</p>
                                </div>
                            )}
                            <video
                                ref={videoRef}
                                width="320"
                                height="240"
                                autoPlay
                                muted
                                playsInline
                                style={{
                                    borderRadius: '8px',
                                    backgroundColor: '#000',
                                    display: cameraReady ? 'block' : 'none'
                                }}
                            />
                            <canvas
                                ref={canvasRef}
                                width="320"
                                height="240"
                                style={{ display: 'none' }}
                            />
                        </div>
                    ) : null}

                    {cameraError && (
                        <div className="w-[320px] p-4 bg-red-100 border border-red-400 text-red-700 rounded">
                            <p>{cameraError}</p>
                            <button
                                onClick={() => {
                                    setCameraError('');
                                    setIsCameraAvailable(true);
                                    setCameraReady(false);
                                }}
                                className="mt-2 text-blue-600 hover:underline"
                            >
                                Thử lại
                            </button>
                        </div>
                    )}
                </div>

                <div>{formattedTime}</div>

                <Form {...form}>
                    <form
                        onSubmit={form.handleSubmit(onSubmit)}
                        className="space-y-4 w-[350px]"
                    >
                        <FormField
                            control={form.control}
                            name="category_id"
                            render={({field}) => (
                                <FormItem>
                                    <FormLabel>Hình thức làm việc</FormLabel>
                                    <Select
                                        onValueChange={field.onChange}
                                        defaultValue={field.value}
                                    >
                                        <FormControl>
                                            <SelectTrigger className="w-full">
                                                <SelectValue placeholder="Chọn hình thức làm việc"/>
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
                                    <FormMessage/>
                                </FormItem>
                            )}
                        />

                        <FormField
                            control={form.control}
                            name="office_id"
                            render={({field}) => (
                                <FormItem>
                                    <FormLabel>Văn phòng</FormLabel>
                                    <Select
                                        onValueChange={field.onChange}
                                        value={field.value}
                                    >
                                        <FormControl>
                                            <SelectTrigger className="w-full">
                                                <SelectValue placeholder="Chọn văn phòng"/>
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
                                    <FormMessage/>
                                </FormItem>
                            )}
                        />

                        <FormField
                            control={form.control}
                            name="note_request"
                            render={({field}) => (
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

                        {/* Hidden fields for gps_location and timestamp */}
                        <FormField
                            control={form.control}
                            name="gps_location"
                            render={({field}) => (
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
                            render={({field}) => (
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
                            render={({field}) => (
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
                            className="w-full mt-2 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
                            disabled={submitAttendanceMutation.isPending}
                        >
                            {submitAttendanceMutation.isPending ? (
                                <span className="flex items-center justify-center">
            <Loading />
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