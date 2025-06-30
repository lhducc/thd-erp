'use client';

import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { useState } from "react";
import { Pencil, Trash2 } from "lucide-react";
import DocumentEmployeeForm from "@/components/DocumentEmployeeForm"

const DocumentDetail = () => {
    const [open, setOpen] = useState(false);

    const Header = () => (
        <div className="flex items-center justify-between bg-white p-4 border-b-2 border-gray-800">
            <div>
                <Label className="text-2xl font-bold">CHI TIẾT TÀI LIỆU</Label>
                <Label className="text-lg font-medium">Nguyễn Văn A</Label>
            </div>
            <Button className="p-3 bg-[#DB290D] text-white rounded-2xl">Xuất file</Button>
        </div>
    );

    const DocumentInfo = () => (
        <div className="bg-white w-600px mt-[20px] rounded-md">
            <p className="text-[25px] py-[10px]">Cam kết bảo mật</p>
            <div className="flex text-[23px]">
                <div className="space-y-3 w-[250px]">
                    <p className="font-semibold  h-[40px] border-b border-gray-500">Mã số:</p>
                    <p className="font-semibold  h-[40px] border-b border-gray-500">Nhóm tài liệu:</p>
                    <p className="font-semibold  h-[40px] border-b border-gray-500">Ngày tạo:</p>
                    <p className="font-semibold  h-[40px] border-b border-gray-500">Hiệu lực từ ngày:</p>
                    <p className="font-semibold  h-[40px] border-b border-gray-500">Ngày hết hạn:</p>
                    <p className="font-semibold  h-[40px] border-b border-gray-500">Tình trạng:</p>
                </div>
                <div className="space-y-3 w-[350px]">
                    <p className=" text-gray-600  h-[40px] border-b border-gray-500">TL000001</p>
                    <p className=" text-gray-600  h-[40px] border-b border-gray-500">Thủ tục tiếp nhận</p>
                    <p className=" text-gray-600  h-[40px] border-b border-gray-500">16/09/2025</p>
                    <p className=" text-gray-600  h-[40px] border-b border-gray-500">16/09/2025</p>
                    <p className=" text-gray-600  h-[40px] border-b border-gray-500">-</p>
                    <p className=" text-gray-600  h-[40px] border-b border-gray-500">Đang hiệu lực</p>
                </div>
            </div>
            <div className="pt-[20px]">
                <p>Nội dung</p>
                <p>Khen thưởng theo quý</p>
            </div>
        </div>
    );


    const DocumentAttachment = () => (
        <div className="mt-8 bg-white p-6 rounded-md mr-[30px]">
            <h3 className="text-xl font-semibold mb-4">Tài liệu đính kèm</h3>
            <div className="border border-gray-300 p-4 rounded-md">
                <p className=" text-gray-600">Tài liệu mẫu PDF sẽ được hiển thị ở đây.</p>
            </div>
        </div>
    );

    return (
        <div className="">
            <Header />
            <div className="flex">
                <DocumentInfo />
                <DocumentAttachment />
            </div>

        </div>
    );
};

export default DocumentDetail;
