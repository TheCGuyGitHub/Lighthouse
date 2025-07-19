<template>
    <nav class="fixed w-full top-0 bg-white shadow px-4 md:px-8">
        <div class="h-14 md:h-16 rounded-md flex items-center justify-between px-4 md:px-6 border my-4">
            <!-- Logo -->
            <NuxtLink to="/" class="flex items-center gap-2 text-2xl font-extrabold text-white tracking-tight">
                <span class="font-mono">Lighthouse</span>
            </NuxtLink>

            <!-- Right Side Icons & Button -->
            <div class="flex items-center gap-4">
                <NuxtLink to="/" class="hover:text-green-400 transition-colors">
                    <User class="w-5 h-5" />
                </NuxtLink>
                <ThemePicker />
            </div>
        </div>
    </nav>
</template>

<script setup lang="ts">
import { User } from 'lucide-vue-next';
import ThemePicker from './ThemePicker.vue';

</script>