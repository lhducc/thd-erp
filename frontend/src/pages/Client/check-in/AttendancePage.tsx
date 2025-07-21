import {useEffect, useRef, useState} from 'react';
import {MapContainer, TileLayer, Marker, Popup} from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';
import {useQuery} from '@tanstack/react-query';
import {getAllOfficesApi} from '@/apis/office.api.ts';
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

const formSchema = z.object({
    work_type: z.enum(['TTS', 'Remote', 'Offline'], {
        required_error: 'Vui lòng chọn hình thức làm việc',
    }),
    office_id: z.string().min(1, 'Vui lòng chọn văn phòng'),
    note: z.string().optional(),
});

const AttendancePage = () => {
    const [position, setPosition] = useState<[number, number] | null>(null);
    const [cameraError, setCameraError] = useState('');
    const [currentTime, setCurrentTime] = useState(new Date());
    const videoRef = useRef<HTMLVideoElement>(null);
    const canvasRef = useRef<HTMLCanvasElement>(null);

    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            work_type: undefined,
            office_id: '',
            note: '',
        },
    });

    const {data: offices = []} = useQuery({
        queryKey: ['offices'],
        queryFn: getAllOfficesApi,
        gcTime: 0,
        staleTime: 0,
    });

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
                },
                (err) => {
                    console.error('Geolocation error:', err);
                }
            );
        }

        if (navigator.mediaDevices?.getUserMedia) {
            navigator.mediaDevices
                .getUserMedia({video: true})
                .then((stream) => {
                    if (videoRef.current) {
                        videoRef.current.srcObject = stream;
                    }
                })
                .catch((err) => {
                    setCameraError('Không thể truy cập camera');
                    console.error('Camera error:', err);
                });
        } else {
            setCameraError('Trình duyệt không hỗ trợ camera');
        }
    }, []);

    useEffect(() => {
        const timer = setInterval(() => {
            setCurrentTime(new Date());
        }, 1000);
        return () => clearInterval(timer);
    }, []);

    const handleTakePhoto = () => {
        if (videoRef.current && canvasRef.current) {
            const ctx = canvasRef.current.getContext('2d');
            if (ctx) {
                ctx.drawImage(videoRef.current, 0, 0, 320, 240);
                const imageData = canvasRef.current.toDataURL('image/png');
                console.log('Ảnh base64:', imageData);
                // TODO: gửi ảnh + form + vị trí về server
            }
        }
    };

    // Format thời gian hiển thị
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

    // Danh sách văn phòng được sắp theo khoảng cách
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
    }, [position, offices]);


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
                    <video
                        ref={videoRef}
                        autoPlay
                        playsInline
                        width="350"
                        height="350"
                        className="rounded shadow"
                    />
                    {cameraError && <p className="text-red-500">{cameraError}</p>}
                    <canvas
                        ref={canvasRef}
                        width="320"
                        height="240"
                        style={{display: 'none'}}
                    />
                </div>

                <div>{formattedTime}</div>

                <Form {...form}>
                    <form
                        onSubmit={form.handleSubmit((data) => {
                            console.log('Form submit:', data);
                            handleTakePhoto();
                        })}
                        className="space-y-4 w-[350px]"
                    >
                        <FormField
                            control={form.control}
                            name="work_type"
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
                                            <SelectItem value="TTS">TTS</SelectItem>
                                            <SelectItem value="Remote">Remote</SelectItem>
                                            <SelectItem value="Offline">Offline</SelectItem>
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
                            name="note"
                            render={({field}) => (
                                <FormItem>
                                    <FormLabel>Ghi chú</FormLabel>
                                    <FormControl>
                                        <input
                                            type="text"
                                            className="input w-full px-3 py-2 border rounded"
                                            {...field}
                                        />
                                    </FormControl>
                                </FormItem>
                            )}
                        />

                        <button
                            type="submit"
                            className="w-full mt-2 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
                        >
                            Chấm công
                        </button>
                    </form>
                </Form>
            </div>
        </div>
    );
};

export default AttendancePage;