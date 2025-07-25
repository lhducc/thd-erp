'use client';

import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { useState } from "react";
import { Pencil, Trash2 } from "lucide-react";
import DocumentEmployeeForm from "@/components/DocumentEmployeeForm";

const DecisionDetail = () => {
    const Header = () => (
        <div className="flex items-center justify-between bg-white p-6 border-b-2 border-gray-800">
            <div>
                <Label className="text-3xl font-semibold">CHI TIẾT QUYẾT ĐỊNH</Label>
                <Label className="text-xl font-medium">Nguyễn Văn A</Label>
            </div>
            <Button className="p-3 bg-[#DB290D] text-white rounded-2xl hover:bg-[#b83a1a]">
                Xuất file
            </Button>
        </div>
    );

    const DocumentInfo = () => (
        <div className="bg-white w-[700px] mt-[10px] rounded-md p-6 ">
            <p className="text-[22px] font-bold py-[10px] border-b border-gray-400">Quyết định về việc ...........................................................</p>
            <div className="flex text-[20px]">
                <div className="space-y-4 w-[300px]">
                    <p className="font-semibold border-b border-gray-400 py-2">Mã số:</p>
                    <p className="font-semibold border-b border-gray-400 py-2">Loại quyết định:</p>
                    <p className="font-semibold border-b border-gray-400 py-2">Nhóm quyết định:</p>
                    <p className="font-semibold border-b border-gray-400 py-2">Ngày tạo:</p>
                    <p className="font-semibold border-b border-gray-400 py-2">Ngày ký:</p>
                    <p className="font-semibold border-b border-gray-400 py-2">Hiệu lực từ ngày:</p>
                    <p className="font-semibold border-b border-gray-400 py-2">Tình trạng:</p>
                    <p className="font-semibold border-b border-gray-400 py-2">Ngày bắt đầu:</p>
                </div>
                <div className="space-y-4 w-[350px]">
                    <p className="text-gray-600 border-b border-gray-400 py-2">QĐ00001</p>
                    <p className="text-gray-600 border-b border-gray-400 py-2">Khen thưởng theo quý</p>
                    <p className="text-gray-600 border-b border-gray-400 py-2">Hình thức khen thưởng</p>
                    <p className="text-gray-600 border-b border-gray-400 py-2">16/09/2025</p>
                    <p className="text-gray-600 border-b border-gray-400 py-2">16/09/2025</p>
                    <p className="text-gray-600 border-b border-gray-400 py-2">16/09/2025</p>
                    <p className="text-gray-600 border-b border-gray-400 py-2">Đang hiệu lực</p>
                    <p className="text-gray-600 border-b border-gray-400 py-2">10/05/2025</p>
                </div>
            </div>
            <div className="pt-[20px]">
                <p className="font-semibold">Nội dung</p>
                <p>Khen thưởng theo quý</p>
            </div>
        </div>
    );

    const DocumentAttachment = () => (
        <div className="mt-8 bg-white p-6 rounded-md mr-[30px] ">
            <h3 className="text-xl font-semibold mb-4">Tài liệu đính kèm</h3>
            <div className="border border-gray-300 p-6 rounded-md">
                <p className="text-gray-600">Tài liệu mẫu PDF sẽ được hiển thị ở đây.</p>
            </div>
        </div>
    );

    return (
        <div className="">
            <Header />
            <div className="flex justify-between">
                <DocumentInfo />
                <DocumentAttachment />
            </div>
        </div>
    );
};

export default DecisionDetail;
