export const InfoRow = ({ label, value }: { label: string; value: string }) => (
    <div className="flex justify-between border-b py-2">
        <span className="font-semibold">{label}</span>
        <span>{value}</span>
    </div>
);