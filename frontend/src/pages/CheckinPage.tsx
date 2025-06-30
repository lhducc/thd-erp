const Checkin = () => {
    return(
        <div className="bg-gray-300 h-full p-[20px] m-0">
            <header>
                <p className="text-4xl font-bold">Chấm công</p>
            </header>
            <div className="w-full h-full flex justify-center mt-[30px]">
                <div className="w-[70%] h-full bg-white rounded-t-3xl justify-center flex">
                    <div>
                        <img src="/image4.png" className="mt-[20px] w-[300px]" alt="" />
                        <div className="w-[300px] bg-[#D9D9D9] h-[300px] mt-[20px] flex items-center justify-center">
                            <img src="/FaceID.png" className="mt-[20px]" alt="" />
                        </div>
                        <p className="text-2xl">Thời gian: 13:22</p>
                        <p className="text-2xl">Thứ 6 ngày 13/06/2025</p>
                        <p>Hình thức chấm công</p>
                        <select className="w-full border text-3xl rounded-2xl">
                            <option>1</option>
                            <option>2</option>
                            <option>3</option>
                            <option>4</option>
                        </select>
                    </div>
                    

                </div>
            </div>

        </div>
    );
}

export default Checkin;