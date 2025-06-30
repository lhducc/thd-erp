import { useState } from "react";

interface ContractFilterProps {
    departments: string[];
    contractTypes: string[];
    conditions: string[];
    onFilterChange: (filters: {
        department: string;
        contractType: string;
        condition: string;
    }) => void;
}

export const ContractFilter = ({
                                   departments,
                                   contractTypes,
                                   conditions,
                                   onFilterChange,
                               }: ContractFilterProps) => {
    const [isOpen, setIsOpen] = useState(false);
    const [selectedFilters, setSelectedFilters] = useState({
        department: "",
        contractType: "",
        condition: "",
    });

    const handleFilterChange = (key: keyof typeof selectedFilters, value: string) => {
        const newFilters = {
            ...selectedFilters,
            [key]: value === selectedFilters[key] ? "" : value, // Toggle selection
        };
        setSelectedFilters(newFilters);
        onFilterChange(newFilters);
    };

    const clearFilters = () => {
        const clearedFilters = {
            department: "",
            contractType: "",
            condition: "",
        };
        setSelectedFilters(clearedFilters);
        onFilterChange(clearedFilters);
    };

    return (
        <div className="relative flex items-center">
            <button
                id="dropdownDefault"
                data-dropdown-toggle="dropdown"
                className="text-white bg-primary-700 hover:bg-primary-800 focus:ring-4 focus:outline-none focus:ring-primary-300 font-medium rounded-lg text-sm px-4 py-2.5 text-center inline-flex items-center dark:bg-primary-600 dark:hover:bg-primary-700 dark:focus:ring-primary-800"
                type="button"
                onClick={() => setIsOpen(!isOpen)}
            >
                <img src="src/assets/Vector.png" alt="filter" />
            </button>

            {/* Dropdown menu */}
            {isOpen && (
                <div
                    id="dropdown"
                    className="z-10 absolute top-12 right-0 w-56 p-3 bg-white rounded-lg shadow dark:bg-gray-700 border"
                >
                    <div className="flex justify-between items-center mb-3">
                        <h6 className="text-sm font-medium text-gray-900 dark:text-white">
                            Bộ lọc hợp đồng
                        </h6>
                        <button
                            onClick={clearFilters}
                            className="text-xs text-primary-600 hover:underline dark:text-primary-400"
                        >
                            Xóa bộ lọc
                        </button>
                    </div>

                    {/* Department filter */}
                    <div className="mb-4">
                        <h6 className="mb-2 text-sm font-medium text-gray-900 dark:text-white">
                            Phòng ban
                        </h6>
                        <ul className="space-y-2 text-sm">
                            {departments.map((dept) => (
                                <li key={dept} className="flex items-center">
                                    <input
                                        id={`dept-${dept}`}
                                        type="checkbox"
                                        checked={selectedFilters.department === dept}
                                        onChange={() => handleFilterChange("department", dept)}
                                        className="w-4 h-4 bg-gray-100 border-gray-300 rounded text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 dark:ring-offset-gray-700 focus:ring-2 dark:bg-gray-600 dark:border-gray-500"
                                    />
                                    <label
                                        htmlFor={`dept-${dept}`}
                                        className="ml-2 text-sm font-medium text-gray-900 dark:text-gray-100"
                                    >
                                        {dept}
                                    </label>
                                </li>
                            ))}
                        </ul>
                    </div>

                    {/* Contract type filter */}
                    <div className="mb-4">
                        <h6 className="mb-2 text-sm font-medium text-gray-900 dark:text-white">
                            Loại hợp đồng
                        </h6>
                        <ul className="space-y-2 text-sm">
                            {contractTypes.map((type) => (
                                <li key={type} className="flex items-center">
                                    <input
                                        id={`type-${type}`}
                                        type="checkbox"
                                        checked={selectedFilters.contractType === type}
                                        onChange={() => handleFilterChange("contractType", type)}
                                        className="w-4 h-4 bg-gray-100 border-gray-300 rounded text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 dark:ring-offset-gray-700 focus:ring-2 dark:bg-gray-600 dark:border-gray-500"
                                    />
                                    <label
                                        htmlFor={`type-${type}`}
                                        className="ml-2 text-sm font-medium text-gray-900 dark:text-gray-100"
                                    >
                                        {type}
                                    </label>
                                </li>
                            ))}
                        </ul>
                    </div>

                    {/* Condition filter */}
                    <div>
                        <h6 className="mb-2 text-sm font-medium text-gray-900 dark:text-white">
                            Tình trạng
                        </h6>
                        <ul className="space-y-2 text-sm">
                            {conditions.map((cond) => (
                                <li key={cond} className="flex items-center">
                                    <input
                                        id={`cond-${cond}`}
                                        type="checkbox"
                                        checked={selectedFilters.condition === cond}
                                        onChange={() => handleFilterChange("condition", cond)}
                                        className="w-4 h-4 bg-gray-100 border-gray-300 rounded text-primary-600 focus:ring-primary-500 dark:focus:ring-primary-600 dark:ring-offset-gray-700 focus:ring-2 dark:bg-gray-600 dark:border-gray-500"
                                    />
                                    <label
                                        htmlFor={`cond-${cond}`}
                                        className="ml-2 text-sm font-medium text-gray-900 dark:text-gray-100"
                                    >
                                        {cond}
                                    </label>
                                </li>
                            ))}
                        </ul>
                    </div>
                </div>
            )}

            {/* Click outside to close dropdown */}
            {isOpen && (
                <div
                    className="fixed inset-0 z-0"
                    onClick={() => setIsOpen(false)}
                />
            )}
        </div>
    );
};