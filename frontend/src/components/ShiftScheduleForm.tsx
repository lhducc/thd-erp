'use client';

import { useState } from 'react';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from '@/components/ui/select';
import { Calendar } from 'lucide-react';

const weekdays = ['T2', 'T3', 'T4', 'T5', 'T6', 'T7', 'CN'];
const weeks = [1, 2, 3, 4, 5];

export default function ShiftScheduleForm() {
  const [schedule, setSchedule] = useState<Record<string, string[]>>({});

  const shiftOptions = [
    'Ca sáng 8:30 - 12:00',
    'Ca chiều 13:30 - 17:30',
    'Full time 8:30 - 17:30',
  ];

  const addShift = (week: number, day: string, shift: string) => {
    const key = `${week}-${day}`;
    setSchedule((prev) => ({
      ...prev,
      [key]: [...(prev[key] || []), shift],
    }));
  };

  const removeShift = (week: number, day: string, index: number) => {
    const key = `${week}-${day}`;
    const updated = [...(schedule[key] || [])];
    updated.splice(index, 1);
    setSchedule((prev) => ({
      ...prev,
      [key]: updated,
    }));
  };

  return (
    <div className="p-6 space-y-6">
      <div className="flex w-full justify-evenly">
        <div className="w-[40%] space-y-2">
          <p>* Tên phân ca lặp</p>
          <Input className="rounded-md" placeholder="* Tên phân ca lặp" />
          <p>* Tình trạng</p>
          <select className="w-full p-2 border border-gray-300 rounded-md text-gray-700 focus:outline-none">
            <option value="">* Tình trạng</option>
            <option value="active">Đang hiệu lực</option>
            <option value="inactive">Chưa hiệu lực</option>
          </select>
          <p>* Ngày hiệu lực</p>
          <div className="relative">
            <Input
              className="rounded-md"
              placeholder="* Ngày hiệu lực"
              type="date"
            />
            <Calendar className="absolute right-3 top-3 w-4 h-4 text-gray-500" />
          </div>
        </div>

        <div className="w-[40%] space-y-2">
          <p>* Đối tượng áp dụng</p>
          <Input className="rounded-md" placeholder="* Đối tượng áp dụng" />
          <p>* Văn phòng</p>
          <Input className="rounded-md" placeholder="* Văn phòng" />
          <p>* Lặp lại theo</p>
          <select className="w-full p-2 border border-gray-300 rounded-md text-gray-700 focus:outline-none">
            <option value="">* Lặp lại theo</option>
            <option value="week">Theo tuần</option>
            <option value="month">Theo tháng</option>
          </select>
        </div>
      </div>

      <div className="overflow-x-auto">
        <table className="table-auto w-full">
          <thead className="rounded-t-2xl bg-red-600">
            <tr className="rounded-t-2xl">
              <th className="w-20 bg-white p-2 text-center"> </th>
              {weekdays.map((day) => (
                <th
                  key={day}
                  className={`p-2 text-center ${
                    day === 'T2' ? 'rounded-tl-2xl' : ''
                  } ${day === 'CN' ? 'rounded-tr-2xl' : ''} text-white`}
                >
                  {day}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {weeks.map((week) => (
              <tr key={week}>
                <td className="p-2 text-center font-semibold">Tuần {week}</td>
                {weekdays.map((day) => {
                  const key = `${week}-${day}`;
                  return (
                    <td
                      key={day}
                      className="border-3 p-2 align-top h-[120px] w-[150px]"
                    >
                      <div className="space-y-2 flex flex-col justify-center items-center h-full">
                        {(schedule[key] || []).map((shift, index) => (
                          <div
                            key={index}
                            className="relative bg-blue-200 text-blue-900 px-3 py-2 rounded-xl max-w-[140px] text-sm"
                          >
                            <div className="font-semibold">Tên ca</div>
                            <div className="mt-1 bg-blue-100 px-2 py-1 rounded-full text-xs text-center">
                              {shift}
                            </div>
                            <button
                              className="absolute top-1.5 right-2 text-sm text-white bg-blue-500 rounded-full w-4 h-4 flex items-center justify-center hover:bg-red-500"
                              onClick={() => removeShift(week, day, index)}
                            >
                              ×
                            </button>
                          </div>
                        ))}
                        <div className="relative">
                          <button
                            onClick={() => setSchedule((prev) => ({ ...prev, [`show-${key}`]: true }))}
                            className={`w-[40px] h-[40px] border border-gray-300 rounded-full text-center text-gray-400 text-xl flex items-center justify-center ${
                              schedule[key]?.length ? 'hidden' : ''
                            }`}
                          >
                            +
                          </button>
                          {schedule[`show-${key}`] && (
                            <div className="absolute top-10 z-10 bg-white border rounded-md shadow-md w-max text-sm">
                              {shiftOptions.map((option) => (
                                <div
                                  key={option}
                                  className="px-4 py-2 hover:bg-blue-100 cursor-pointer"
                                  onClick={() => {
                                    addShift(week, day, option);
                                    setSchedule((prev) => {
                                      const updated = { ...prev };
                                      delete updated[`show-${key}`];
                                      return updated;
                                    });
                                  }}
                                >
                                  {option}
                                </div>
                              ))}
                            </div>
                          )}
                        </div>


                      </div>
                    </td>
                  );
                })}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
