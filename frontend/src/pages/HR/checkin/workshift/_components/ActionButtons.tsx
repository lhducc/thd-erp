import WorkshiftForm from "@/pages/HR/checkin/workshift/_components/WorkshiftForm.tsx";
import Loading from "@/components/Loading.tsx";
import {Suspense} from "react";
import {Button} from "@/components/ui/button.tsx";
import {EyeIcon, SquarePen} from "lucide-react";
import type {Workshift} from "@/types/workshift.ts";
import ConfirmDelete from "@/components/ConfirmDelete.tsx";

export const ActionButtons = ({workshift, refetchWorkshifts, deleteWorkshift}: {
    workshift: Workshift;
    refetchWorkshifts: () => void;
    deleteWorkshift: (id: string) => void;
}) => {
    return (
        <div className="flex gap-4">
            <Suspense fallback={<Loading/>}>
                <WorkshiftForm
                    editBtn={
                        <Button variant="outline">
                            <EyeIcon/>
                        </Button>
                    }
                    data={workshift}
                    type="view"
                    refetch={refetchWorkshifts}
                />
            </Suspense>
            <Suspense fallback={<Loading/>}>
                <WorkshiftForm
                    editBtn={
                        <Button variant="outline">
                            <SquarePen/>
                        </Button>
                    }
                    data={workshift}
                    type="edit"
                    refetch={refetchWorkshifts}
                />
            </Suspense>
            <ConfirmDelete deleteFn={() => deleteWorkshift(workshift.workshift_id)}/>
        </div>
    );
};