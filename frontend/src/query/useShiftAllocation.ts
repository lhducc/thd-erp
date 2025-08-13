import {useQuery} from "@tanstack/react-query";
import {getAllShiftAllocations} from "@/apis/shift-allocation.api.ts";
import {useState} from "react";

export const useGetAllShiftAllocations = (month: number, year: number) => {
    const [page] = useState(1);
    const [pageSize] = useState(100);

    return useQuery({
        queryKey: ["shift-allocation", month, year],
        queryFn: () => getAllShiftAllocations(month, year, page, pageSize),
    });
};