import {useState} from "react";
import {Button} from "@/components/ui/button";
import {SquarePen} from "lucide-react";
import ConfirmDelete from "@/components/ConfirmDelete";
import {ContractForm} from "@/components/ContractForm";

type Props = {
    contractId: string;
};

const deleteContract = async (id: string) => {
    console.log("Delete contract", id);
};

export const ContractActions = ({contractId}: Props) => {
    const [open, setOpen] = useState(false);

    return (
        <div className="flex gap-2">
            <Button variant="outline" onClick={() => setOpen(!open)}>
                <SquarePen />
            </Button>
            <ContractForm open={open} setOpen={setOpen} contractId={contractId} />
            <ConfirmDelete deleteFn={() => deleteContract(contractId)} />
        </div>
    );
};
