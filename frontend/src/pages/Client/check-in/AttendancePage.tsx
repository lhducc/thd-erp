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
import {toVietnamISOString} from "@/lib/utils.ts";

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
    const [position, setPosition] = useState<[number, number] | null>(null);

    const [currentTime, setCurrentTime] = useState(new Date());
    const videoRef = useRef<HTMLVideoElement>(null);
    const canvasRef = useRef<HTMLCanvasElement>(null);

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

    const {data: offices = [], isLoading: isOfficesLoading} = useOffice();

    const startCamera = async () => {
        try {
            const stream = await navigator.mediaDevices.getUserMedia({ video: { facingMode: 'user' } });
            console.log("stream", stream);
            if (videoRef.current) {
                console.log('videoref')
                videoRef.current.srcObject = stream;
                console.log(videoRef.current);
            }
            return true;
        } catch (err) {
            console.error("Không thể mở camera:", err);
            toast.error("Không thể mở camera");
            return false;
        }
    };

    const stopCamera = () => {
        const stream = videoRef.current?.srcObject as MediaStream;
        stream?.getTracks().forEach(track => track.stop());
    };

    useEffect(() => {
        if (typeof navigator !== "undefined" && navigator.mediaDevices?.getUserMedia) {
            let stream: MediaStream | null = null;
            const initCamera = async () => {
                try {
                    stream = await navigator.mediaDevices.getUserMedia({
                        video: { facingMode: 'user' }
                    });
                    if (videoRef.current) {
                        videoRef.current.srcObject = stream;
                    }
                } catch (err) {
                    console.error("Camera error:", err);
                    toast.error("Không thể mở camera");
                }
            };
            initCamera();

            return () => {
                stream?.getTracks().forEach(track => track.stop());
            };
        } else {
            console.error("Camera API không khả dụng");
            toast.error("Trình duyệt không hỗ trợ camera hoặc cần HTTPS");
        }
    }, []);


    const handleSubmit = async () => {
        if (!canvasRef.current || !videoRef.current) return;
        const context = canvasRef.current.getContext("2d");
        if (!context) return;

        context.drawImage(videoRef.current, 0, 0, canvasRef.current.width, canvasRef.current.height);

        canvasRef.current.toBlob(async (blob) => {
            if (blob) {
                const file = new File([blob], `attendance_${Date.now()}.jpg`, { type: "image/jpeg" });
                form.setValue("image", file);

                // submit form
                await form.handleSubmit(onSubmit)();
            }
        }, "image/jpeg");
    };


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
        },
        onError: (error) => {
            if (axios.isAxiosError(error)) {
                console.log(error.response);
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
            console.log("hi")
            console.log(submissionData)
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

    if (isCategoriesLoading || isOfficesLoading) {
        return <Loading />;
    }

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

                <video
                    ref={videoRef}
                    autoPlay
                    muted
                    playsInline
                    className="rounded-lg border border-gray-300"
                    width="300"
                    height="300"
                />
                <canvas ref={canvasRef} width={300} height={300} className="hidden" />

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