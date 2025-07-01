"use client";

import * as React from "react";
import {
    type ColumnDef,
    type ColumnFiltersState,
    type SortingState,
    type VisibilityState,
    getCoreRowModel,
    getFilteredRowModel,
    getPaginationRowModel,
    getSortedRowModel,
    useReactTable,
} from "@tanstack/react-table";

import {InputWithIcon} from "@/components/ui/input-with-icon.tsx";

interface DataTableProps<TData, TValue> {
    columns: ColumnDef<TData, TValue>[];
    data: TData[];
}

export default function DataTable<TData, TValue>({
        columns,
        data,
        title,
        buttonCreate,
        navLink,
        keyFilter,
    }: DataTableProps<TData, TValue> & {
    title?: string;
    buttonCreate?: React.ReactNode;
    navLink?: React.ReactNode;
    keyFilter?: string;
}) {
    const [sorting, setSorting] = React.useState<SortingState>([]);
    const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
        []
    );
    const [columnVisibility, setColumnVisibility] =
        React.useState<VisibilityState>({});
    const [rowSelection, setRowSelection] = React.useState({});
    const [pagination, setPagination] = React.useState({
        pageIndex: 0,
        pageSize: 10,
    });
    console.log("datdataa",data);
    const table = useReactTable({
        data,
        columns,
        onSortingChange: setSorting,
        onColumnFiltersChange: setColumnFilters,
        getCoreRowModel: getCoreRowModel(),
        getPaginationRowModel: getPaginationRowModel(),
        getSortedRowModel: getSortedRowModel(),
        getFilteredRowModel: getFilteredRowModel(),
        onColumnVisibilityChange: setColumnVisibility,
        onRowSelectionChange: setRowSelection,
        onPaginationChange: setPagination,
        state: {
            sorting,
            columnFilters,
            columnVisibility,
            rowSelection,
            pagination,
        },
    });

    return (
        <div className="w-full">
            <div className="flex lg:flex-row flex-col items-center py-4 gap-8">
                <div className="lg:text-3xl text-lg font-bold text-nowrap uppercase">{title}</div>
                <InputWithIcon
                    value={
                        (table.getColumn(keyFilter || "")?.getFilterValue() as string) ?? ""
                    }
                    onChange={(event) =>
                        table.getColumn(keyFilter || "")?.setFilterValue(event.target.value)
                    }
                    className="w-full h-[50px]"
                />
                {buttonCreate}
            </div>
            {navLink}
        </div>
    );
}
