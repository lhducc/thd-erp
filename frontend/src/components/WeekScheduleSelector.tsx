import { Checkbox } from "@/components/ui/checkbox"
import { Label } from "@/components/ui/label"
import { Card, CardHeader, CardContent, CardTitle, CardDescription } from "@/components/ui/card"
import { Clock, CalendarDays } from "lucide-react"
import { Badge } from "@/components/ui/badge"

export const WorkshiftScheduler = ({ workshifts }) => {
    const days = [
        { id: 'mon', name: 'Thứ 2' },
        { id: 'tue', name: 'Thứ 3' },
        { id: 'wed', name: 'Thứ 4' },
        { id: 'thu', name: 'Thứ 5' },
        { id: 'fri', name: 'Thứ 6' },
        { id: 'sat', name: 'Thứ 7' },
        { id: 'sun', name: 'Chủ nhật' }
    ]

    const formatTime = (timeStr) => {
        return timeStr.substring(0, 5) // Converts "HH:mm:ss" to "HH:mm"
    }

    return (
        <div className="flex flex-col lg:flex-row gap-6 p-4 bg-gray-50 rounded-lg">
            {/* Days Column - Fixed height with overflow */}
            <Card className="lg:w-1/4">
                <CardHeader>
                    <div className="flex items-center gap-2">
                        <CalendarDays className="w-5 h-5 text-blue-600" />
                        <CardTitle>Ngày làm việc</CardTitle>
                    </div>
                    <CardDescription>Chọn ngày cần áp dụng lịch</CardDescription>
                </CardHeader>
                <CardContent className="overflow-y-auto">
                    <div className="space-y-3 p-1">
                        {days.map(day => (
                            <div key={day.id} className="flex items-center space-x-3 p-3 hover:bg-gray-100 rounded-lg transition-colors">
                                <Checkbox id={day.id} />
                                <Label htmlFor={day.id} className="text-base cursor-pointer">
                                    {day.name}
                                </Label>
                            </div>
                        ))}
                    </div>
                </CardContent>
            </Card>

            {/* Workshifts Column - Fixed height with overflow */}
            <div className="lg:w-3/4 space-y-4">
                <Card>
                    <CardHeader>
                        <div className="flex items-center gap-2">
                            <Clock className="w-5 h-5 text-blue-600" />
                            <CardTitle>Ca làm việc</CardTitle>
                        </div>
                        <CardDescription>Chọn ca làm việc cho ngày đã chọn</CardDescription>
                    </CardHeader>
                </Card>

                <div className="overflow-y-auto border rounded-lg" style={{ maxHeight: '400px' }}>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4 p-4">
                        {workshifts?.map(shift => (
                            <Card key={shift.workshift_id} className="hover:shadow-md transition-shadow">
                                <CardContent className="p-4">
                                    <div className="flex items-start space-x-3">
                                        <Checkbox id={`shift-${shift.workshift_id}`} className="mt-1" />
                                        <div className="flex-1">
                                            <div className="flex items-center justify-between">
                                                <Label htmlFor={`shift-${shift.workshift_id}`} className="text-base font-medium">
                                                    {shift.workshift_name}
                                                </Label>
                                                <Badge variant="outline">{shift.time_of_day}</Badge>
                                            </div>
                                            <div className="mt-2 space-y-1 text-sm text-gray-600">
                                                <div className="flex">
                                                    <span className="w-24">Thời gian:</span>
                                                    <span>
                            {formatTime(shift.start_time)} - {formatTime(shift.end_time)}
                          </span>
                                                </div>
                                                <div className="flex">
                                                    <span className="w-24">Check-in:</span>
                                                    <span>
                            {formatTime(shift.checkin_from)} - {formatTime(shift.checkin_to)}
                          </span>
                                                </div>
                                                <div className="flex">
                                                    <span className="w-24">Check-out:</span>
                                                    <span>
                            {formatTime(shift.checkout_from)} - {formatTime(shift.checkout_to)}
                          </span>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>
                        ))}
                    </div>
                </div>

                {/*<div className="flex justify-end gap-3 pt-2">*/}
                {/*    <Button variant="outline">Hủy bỏ</Button>*/}
                {/*    <Button>Xác nhận lịch</Button>*/}
                {/*</div>*/}
            </div>
        </div>
    )
}

export default WorkshiftScheduler