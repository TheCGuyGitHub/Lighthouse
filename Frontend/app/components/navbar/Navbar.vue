<template>
    <nav class="fixed w-full shadow px-4 md:px-8 bg-background">
        <div class="h-14 md:h-16 rounded-md flex items-center justify-between px-4 md:px-6 border my-4 border-border">
            <!-- Logo -->
            <NuxtLink to="/" class="flex items-center gap-2 text-2xl font-extrabold text-foreground tracking-tight">
                <span class="font-mono">Lighthouse</span>
            </NuxtLink>

            <!-- Right Side Icons & Button -->
            <div class="flex items-center gap-2">
                <NuxtLink to="/pages" class="flex items-center gap-1 text-foreground hover:text-primary">
                    <BookOpenText class="w-5 h-5" />
                </NuxtLink>
                <NuxtLink to="/servers" class="flex items-center gap-1 text-foreground hover:text-primary">
                    <Server class="w-5 h-5" />
                </NuxtLink>
                <Profile />
                <ThemePicker />
            </div>
        </div>
    </nav>
</template>

<script setup lang="ts">
import ThemePicker from './ThemePicker.vue';
import Profile from './Profile.vue';
import { BookOpenText, Server } from 'lucide-vue-next';

</script>