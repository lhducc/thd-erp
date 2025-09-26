import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button.tsx";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import InputPassword from "@/components/InputPassword.tsx";
import LogoSignInPage from "@/assets/LogoSignInPage.svg";
import { useMutation } from "@tanstack/react-query";
import { signInApi } from "@/apis/signIn.api.ts";
import { toast } from "sonner";
import { useNavigate } from "react-router-dom";
import { useEffect } from "react";
import Loading from "@/components/Loading";
import PATH from "@/constants/Path";

const formSchema = z.object({
  email: z.string().email("Email không hợp lệ."),
  password: z.string().nonempty("Vui lòng nhập mật khẩu."),
});

const SignInPage = () => {

  const navigate = useNavigate(); 

  useEffect(() => {
    const token = localStorage.getItem("access_token");
    if (token) {
      navigate("/");
    }
  }, [navigate]);

  // 1. Define your form.
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const { mutate: login, isPending } = useMutation({
    mutationFn: signInApi,
    onSuccess: (res) => {
      if (res.statuscode === 200) {
        localStorage.setItem("access_token", res.data.access_token);
        localStorage.setItem("refresh_token", res.data.refresh_token);
        
        // Xử lý chuyển hướng dựa trên first_login
        if (res.data.first_login) {
          toast.success("Đăng nhập thành công. Vui lòng đổi mật khẩu lần đầu.");
          navigate(PATH.FIRST_CHANGE_PASSWORD); // Đường dẫn đến trang đổi mật khẩu đầu tiên
        } else {
          toast.success("Đăng nhập thành công");
          navigate("/"); 
        }
      }
    },
    onError: (error: any) => {
      toast.error(error.message || "Đăng nhập thất bại");
    },
  });

  // 2. Define a submit handler.
  function onSubmit(values: z.infer<typeof formSchema>) {
    login(values)
  }

  return (
    <main className="bg-[#FFF2F0] min-h-screen flex justify-center items-center md:gap-20 md:flex-row flex-col">
      <img src={LogoSignInPage} alt="logo" className="max-md:w-[250px]" />
      <Card className="border-[#B50101] border-2 md:w-[400px] w-[90%]">
        <CardHeader>
          <CardTitle className="text-center text-[#B50101]">
            Đăng nhập
          </CardTitle>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="example@thdcybersecurity.xyz"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Mật khẩu</FormLabel>
                    <FormControl>
                      <InputPassword field={field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button type="submit" className={`w-full`} disabled={isPending}>
                {!isPending ? "Đăng nhập" : <Loading />}
              </Button>
            </form>
          </Form>
        </CardContent>
      </Card>
    </main>
  );
};

export default SignInPage;
