import { Loader2 } from "lucide-react";

const Loading = () => {
  return (
    <div className="flex h-full flex-1 items-center justify-center">
      <Loader2 className="h-10 w-10 animate-spin" />
    </div>
  );
};

export default Loading;
