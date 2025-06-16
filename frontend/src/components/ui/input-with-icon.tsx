import { Input } from "@/components/ui/input";
import searchIcon from "@/assets/seach.svg";

export function InputWithIcon({
                                  value,
                                  onChange,
                                  placeholder = "Tìm kiếm...",
                                  className,
                              }: {
    value: string;
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
    placeholder?: string;
    className?: string;
}) {
    return (
        <div className={`relative ${className}`}>
            <div className="absolute left-3 top-1/2 -translate-y-1/2">
                <img src={searchIcon} alt="Search icon" width={16} height={16} />
            </div>
            <Input
                className="pl-10 h-[50px]"
                placeholder={placeholder}
                value={value}
                onChange={onChange}
            />
        </div>
    );
}