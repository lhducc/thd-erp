import {useQuery} from "@tanstack/react-query";
import {getJobTitles} from "@/apis/jobTitle.api.ts";

export const useJobTitle = () =>
    useQuery({
    queryKey: ["jobTitles"],
    queryFn: getJobTitles,
});