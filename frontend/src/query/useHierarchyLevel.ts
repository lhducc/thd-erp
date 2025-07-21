import {useQuery} from "@tanstack/react-query";
import {getAllHierarchyLevelApi} from "@/apis/hierarchyLevel.api.ts";

export const useHierarchyLevel = () =>
    useQuery({
        queryKey: ["hierarchy-level"],
        queryFn: getAllHierarchyLevelApi,
    });