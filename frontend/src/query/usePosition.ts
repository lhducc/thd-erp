import {useQuery} from "@tanstack/react-query";
import {getAllPositionApi} from "@/apis/position.api.ts";

export const usePosition = () =>
    useQuery({
        queryKey: ["positions"],
        queryFn: getAllPositionApi,
    });