'use client';
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { useState } from "react";
import { Pencil, Trash2 } from "lucide-react";

const InsuranceInformationPage = () => {
    const [open, setOpen] = useState(false);

    const Header = () => (
        <div className="flex items-center justify-between bg-white p-4 border-b-1 border-gray-800">
            <Label className="text-2xl font-bold">THÔNG TIN BẢO HIỂM</Label>
            <Button className="bg-[#DB290D] hover:bg-[red] text-white rounded-lg text-xl px-4 py-2">
                Chỉnh sửa thông tin
            </Button>
        </div>
    );

    const GeneralInfo = () => (
        <div className="bg-white p-6 rounded-md w-[45%]">
            <h2 className="text-xl font-semibold mb-4">Thông tin chung</h2>
            <div className="space-y-2">
                <div className="flex justify-between w-full border-b-1 border-gray-400 py-2">
                    <p className=" text-xl">Mức đóng</p>
                    <p className="font-medium text-left w-1/2">xxxxxxxxxxxxxxxx</p>
                </div>
                <div className="flex justify-between w-full border-b-1 border-gray-400 py-2">
                    <p className=" text-xl">Thời gian bắt đầu</p>
                    <p className="font-medium text-left w-1/2">10/05/2025</p>
                </div>
            </div>
        </div>
    );

    const DataTable = ({ title, columns, data }: any) => (
        <div className="bg-white p-4 rounded-md w-full">
            <h2 className="text-lg font-bold mb-2">{title}</h2>
            <table className="w-full border text-sm">
                <thead className="bg-gray-100">
                    <tr>
                        {columns.map((col: string, i: number) => (
                            <th key={i} className="border px-3 py-2 text-left">{col}</th>
                        ))}
                        <th className="border px-3 py-2">Chỉnh sửa</th>
                    </tr>
                </thead>
                <tbody>
                    {data.map((row: any, idx: number) => (
                        <tr key={idx}>
                            {Object.values(row).map((val: any, i: number) => (
                                <td key={i} className="border px-3 py-2">{val}</td>
                            ))}
                            <td className="border px-3 py-2 text-center">
                                <div className="flex justify-center gap-2">
                                    <Pencil className="w-4 h-4 cursor-pointer" />
                                    <Trash2 className="w-4 h-4 cursor-pointer text-red-500" />
                                </div>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
            <div className="flex justify-center mt-3 space-x-1">
                <button className="px-2 py-1 border rounded hover:bg-gray-200">
                    {"<"}
                </button>
                {[1, 2, 3, 4, 5].map((page) => (
                    <button key={page} className="px-2 py-1 border rounded hover:bg-gray-200">
                        {page}
                    </button>
                ))}
                <button className="px-2 py-1 border rounded hover:bg-gray-200">
                    {">"}
                </button>
            </div>
        </div>
    );

    const companyData = [
        { label: "Bảo hiểm xã hội", value: "17.5" },
        { label: "Bảo hiểm TNLD - BNN", value: "0.5" },
        { label: "Bảo hiểm y tế", value: "3" },
        { label: "Bảo hiểm thất nghiệp", value: "1" }
    ];

    const employeeData = [
        { label: "Bảo hiểm xã hội", value: "8" },
        { label: "Bảo hiểm TNLD - BNN", value: "0" },
        { label: "Bảo hiểm y tế", value: "1.5" },
        { label: "Bảo hiểm thất nghiệp", value: "1" }
    ];

    return (
        <div className="space-y-4">
            <Header />
            <div className="">
                <GeneralInfo />
                <div className="flex">
                    <div className="bg-white p-6 rounded-md w-[45%]">
                        <h2 className="text-xl text-white font-semibold mb-4 rounded-t-lg bg-[#DB290D] py-2 px-[30px]">Công ty chịu thuế</h2>
                        <div className="space-y-2">
                            {companyData.map((item, idx) => (
                                <div key={idx} className="flex justify-between w-full border-b-1 border-gray-400 py-2 px-[30px]">
                                    <p className=" text-sm">{item.label}</p>
                                    <p className="font-medium text-left">{item.value}%</p>
                                </div>
                            ))}
                            <div className="flex justify-between w-full py-2 font-semibold text-lg">
                                <p className="  px-[30px]">Tổng cộng</p>
                                <p className="px-[30px]">{companyData.reduce((acc, curr) => acc + parseFloat(curr.value), 0)}%</p>
                            </div>
                        </div>
                    </div>

                    <div className="bg-white p-6 rounded-md w-[45%]">
                        <h2 className="text-xl text-white font-semibold mb-4 rounded-t-lg bg-[#DB290D] py-2 px-[30px]">Nhân viên chịu thuế</h2>
                        <div className="space-y-2">
                            {employeeData.map((item, idx) => (
                                <div key={idx} className="flex justify-between w-full border-b-1 border-gray-400 py-2 px-[30px]">
                                    <p className=" text-sm">{item.label}</p>
                                    <p className="font-medium text-left">{item.value}%</p>
                                </div>
                            ))}
                            <div className="flex justify-between w-full py-2 font-semibold text-lg">
                                <p className=" px-[30px]">Tổng cộng</p>
                                <p className="px-[30px]">{employeeData.reduce((acc, curr) => acc + parseFloat(curr.value), 0)}%</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default InsuranceInformationPage;
