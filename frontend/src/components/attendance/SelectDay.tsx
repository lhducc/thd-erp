const SelectDay = () => {
    return (
        <div className="w-full h-full justify-center flex cursor-pointer">
            <div className={`bg-green-500/80 rounded-full w-7 h-7 md:w-10 md:h-10 flex justify-center`}>
                <p className={`text-white font-semibold md:text-3xl text-xl`}>+</p>
            </div>
        </div>
    );
};

export default SelectDay;
