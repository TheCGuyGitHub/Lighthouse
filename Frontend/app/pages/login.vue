<template>
    <div class="flex flex-col items-center justify-center h-screen">
        <Card class="mx-auto max-w-sm">
            <CardHeader>
                <CardTitle class="text-2xl">
                    Login
                </CardTitle>
                <CardDescription>
                    Enter your email below to login to your account
                </CardDescription>
            </CardHeader>
            <CardContent>
                <div class="grid gap-4">
                    <div class="grid gap-2">
                        <Label for="email">Email</Label>
                        <Input id="email" type="email" placeholder="m@example.com" required v-model="email" />
                    </div>
                    <div class="grid gap-2">
                        <div class="flex items-center">
                            <Label for="password">Password</Label>
                            <a href="#" class="ml-auto inline-block text-sm underline">
                                Forgot your password?
                            </a>
                        </div>
                        <Input id="password" type="password" required v-model="password" />
                    </div>
                    <Button type="submit" class="w-full" @click="login">
                        Login
                    </Button>
                    <Button variant="outline" class="w-full">
                        Login with Google
                    </Button>
                </div>
                <div class="mt-4 text-center text-sm">
                    Don't have an account?
                    <a href="#" class="underline">
                        Sign up
                    </a>
                </div>
            </CardContent>
        </Card>

    </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Card } from '@/components/ui/card';
import CardHeader from '~/components/ui/card/CardHeader.vue';
import pb from '~/lib/pocketbase';
definePageMeta({
    layout: 'login',
});

const email = ref('')
const password = ref('')
const showPassword = ref(false)
const router = useRouter();

const login = () => {
    pb.collection('users').authWithPassword(email.value, password.value)
        .then((res) => {
            router.push('/');
        })
        .catch((err) => {
            console.error('Login failed:', err);
            // Handle login error
        });
};
</script>