import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export default function Home() {
  return (
    <div className="flex justify-center p-6">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>My Fitness</CardTitle>
          <CardDescription>
            Release 0 — окружение готово: shadcn/ui, TanStack Query, тёмная тема.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Button>Начать</Button>
        </CardContent>
      </Card>
    </div>
  );
}
