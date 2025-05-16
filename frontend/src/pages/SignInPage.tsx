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
import { Input } from "@/components/ui/input";
import InputPassword from "@/components/InputPassword";
import LogoSignInPage from "@/assets/LogoSignInPage.svg";

const formSchema = z.object({
  email: z.string().email("Email không hợp lệ."),
  password: z.string().nonempty("Vui lòng nhập mật khẩu."),
});

const SignInPage = () => {
  // 1. Define your form.
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  // 2. Define a submit handler.
  function onSubmit(values: z.infer<typeof formSchema>) {
    // Do something with the form values.
    // ✅ This will be type-safe and validated.
    console.log(values);
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
              <Button type="submit" className="w-full">
                Đăng nhập
              </Button>
            </form>
          </Form>
        </CardContent>
      </Card>
    </main>
  );
};

export default SignInPage;
