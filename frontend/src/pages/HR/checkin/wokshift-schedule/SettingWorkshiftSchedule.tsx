import {InfoRow} from "@/components/ui/info-row.tsx";

const SettingWorkshiftSchedule = () => {
    return (
        <div>
            <div className="container mx-auto py-6">
                <p className={`text-3xl font-semibold text-gray-800 mb-5`}>Cấu hình lịch làm việc</p>
                <div className="border-t-2 border-gray-500 flex justify-between gap-5">
                    <div className="w-full">
                        <p className={`text-xl font-semibold text-gray-600`}>Lịch làm việc</p>
                        <InfoRow label="Văn phòng" value={`Hồ Bá Kiện`}/>
                    </div>
                    <div className="w-full">
                        <InfoRow label="Ngày hiệu lực" value={'01-07-2025'}/>
                        <InfoRow label="Ngày hết hiệu lực" value={'01-07-2025'}/>
                        <InfoRow label="Tình trạng" value={`Đang hiệu lực`}/>
                    </div>
                </div>
                <div>
                    <p className={`font-semibold`}>Quản lý lịch làm việc</p>
                </div>
            </div>
        </div>
    );
};

export default SettingWorkshiftSchedule;
