import MonthYearPicker from "@/components/MonthYearPicker.tsx";
import {Button} from "@/components/ui/button.tsx";
import {Menu} from "lucide-react";

export const SecondNav = ({
                       currentDate,
                       scheduleNames,
                       selectedSchedule,
                       onDateChange,
                       onToggleSidebar,
                       onSelectSchedule,
                       isSidebarOpen
                   }: {
    currentDate: Date;
    scheduleNames: string[];
    selectedSchedule: string | null;
    onDateChange: ({ month, year }: { month: number; year: number }) => void;
    onToggleSidebar: () => void;
    onSelectSchedule: (name: string) => void;
    isSidebarOpen: boolean;
}) => (
    <div className="flex justify-between items-center py-2">
        <MonthYearPicker
            initialMonth={currentDate.getMonth() + 1}
            initialYear={currentDate.getFullYear()}
            onChange={onDateChange}
        />
        <div className="flex gap-2">
            {scheduleNames.map((name, index) => (
                <button
                    key={index}
                    className={`rounded-sm px-3 py-1 ${
                        selectedSchedule === name
                            ? "bg-blue-500 text-white"
                            : "bg-green-100 hover:bg-green-200"
                    }`}
                    onClick={() => onSelectSchedule(name)}
                >
                    {name}
                </button>
            ))}
        </div>
        <Button variant="ghost" size="icon" onClick={onToggleSidebar}>
            <Menu className="h-5 w-5" />
        </Button>
    </div>
);
