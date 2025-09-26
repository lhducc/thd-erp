import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import InputPassword from "@/components/InputPassword";
import LogoSignInPage from "@/assets/LogoSignInPage.svg";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { useNavigate } from "react-router-dom";
import Loading from "@/components/Loading";
import { changePasswordApi } from "@/apis/signIn.api";
import { useEffect } from "react";

const formSchema = z
  .object({
    current_password: z.string().nonempty("Vui lòng nhập mật khẩu hiện tại."),
    new_password: z.string().nonempty("Vui lòng nhập mật khẩu mới."),
    retype_password: z.string().nonempty("Vui lòng nhập lại mật khẩu mới."),
  })
  .refine((data) => data.new_password === data.retype_password, {
    message: "Mật khẩu mới không khớp.",
    path: ["retype_password"],
  });

const FirstChangePasswordPage = () => {
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem("access_token");
    if (!token) {
      navigate("/sign-in");
    }
  }, [navigate]);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      current_password: "",
      new_password: "",
      retype_password: "",
    },
  });

  const { mutate: changePassword, isPending } = useMutation({
  mutationFn: changePasswordApi,
  onSuccess: (res) => {
    console.log("🔧 DEBUG: API Response:", res);

    // Nếu API chuẩn BaseResponse { statuscode, message, data }
    if (res.statuscode === 200) {
      toast.success(res.message || "Đổi mật khẩu thành công");

      // clear token để buộc user login lại
      localStorage.removeItem("access_token");
      localStorage.removeItem("refresh_token");

      setTimeout(() => {
        navigate("/sign-in", { replace: true }); // quay về login
      }, 1500);
    } else {
      toast.error(res.message || "Mật khẩu cũ không đúng, vui lòng thử lại");
    }
  },
  onError: (error: any) => {
    if (error.response?.status === 401) {
      toast.error("Mật khẩu cũ không đúng, vui lòng thử lại");
      return; 
    }

    const errorMessage =
      error.response?.data?.message || error.message || "Đổi mật khẩu thất bại";
    toast.error(errorMessage);
  },
});

  function onSubmit(values: z.infer<typeof formSchema>) {
    console.log("🔧 DEBUG: Form values:", values); 
    changePassword({
      current_password: values.current_password,
      new_password: values.new_password,
    });
  }

  return (
    <main className="bg-[#FFF2F0] min-h-screen flex justify-center items-center md:gap-20 md:flex-row flex-col">
      <img src={LogoSignInPage} alt="logo" className="max-md:w-[250px]" />
      <Card className="border-[#B50101] border-2 md:w-[400px] w-[90%]">
        <CardHeader>
          <CardTitle className="text-center text-[#B50101]">
            Đổi Mật Khẩu Lần Đầu
          </CardTitle>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
              <FormField
                control={form.control}
                name="current_password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Mật khẩu hiện tại</FormLabel>
                    <FormControl>
                      <InputPassword field={field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="new_password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Mật khẩu mới</FormLabel>
                    <FormControl>
                      <InputPassword field={field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="retype_password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Nhập lại mật khẩu mới</FormLabel>
                    <FormControl>
                      <InputPassword field={field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button type="submit" className="w-full" disabled={isPending}>
                {isPending ? <Loading /> : "Thay đổi"}
              </Button>
            </form>
          </Form>
        </CardContent>
      </Card>
    </main>
  );
};

export default FirstChangePasswordPage;