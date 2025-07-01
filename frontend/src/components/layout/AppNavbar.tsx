import bell from "@/assets/bell.svg"
import avatar from "@/assets/user-avatar.svg"
import {useAuth} from "@/context/AuthContext.tsx";
import {useEffect, useRef, useState} from "react";
import {SidebarTrigger} from "@/components/ui/sidebar.tsx";

const AppNavbar = () => {
    const {currentUser, handleLogout} = useAuth()
    const [isOpen, setIsOpen] = useState(false)
    const dropdownRef = useRef<HTMLDivElement>(null)

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
                setIsOpen(false)
            }
        }

        document.addEventListener("mousedown", handleClickOutside)
        return () => {
            document.removeEventListener("mousedown", handleClickOutside)
        }
    }, [])

    return (
        <div className="h-[80px] flex items-center justify-between gap-5 w-full md:justify-end p-3">
            <SidebarTrigger className={`md:hidden`} />
            <div className="flex items-center justify-between gap-5 md:gap-16">
                <button className="flex items-center justify-between border-gray-300 rounded-sm border-2 p-2 h-fit">
                    <img src={bell} className="md:w-[25px] md:h-[25px] w-[20px] h-[20px]" alt="" />
                </button>
                <div className={`flex gap-5`}>
                    <div className="flex flex-col items-end justify-end">
                        <p className="md:text-[20px] font-semibold w-[150px] md:w-full truncate">{currentUser?.full_name}</p>
                        <p className="md:text-[14px] text-[10px] font-semibold">{currentUser?.roles}</p>
                    </div>
                    <div className="relative" ref={dropdownRef}>
                        <button
                            onClick={() => setIsOpen(!isOpen)}
                            className="rounded-full bg-gray-400 md:w-[50px] md:h-[50px] w-[40px] h-[40px] p-2"
                        >
                            <img src={avatar} alt="avatar" />
                        </button>
                        <div
                            className={`flex bg-white border shadow min-h-10 p-3 right-0 top-16 absolute rounded-lg min-w-[200px] items-center justify-between gap-5 ${
                                isOpen ? "" : "hidden"
                            }`}
                        >
                            <button
                                onClick={handleLogout}
                                className="hover:bg-gray-100 w-full h-full py-2 rounded-lg"
                            >
                                Đăng xuất
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default AppNavbar;
