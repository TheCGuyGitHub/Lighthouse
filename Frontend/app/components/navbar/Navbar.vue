<template>
    <nav class="fixed w-full top-0 shadow px-4 md:px-8">
        <div class="h-14 md:h-16 rounded-md flex items-center justify-between px-4 md:px-6 border my-4 border-gray-900 dark:border-gray-700">
            <!-- Logo -->
            <NuxtLink to="/" class="flex items-center gap-2 text-2xl font-extrabold text-black dark:text-white tracking-tight">
                <span class="font-mono">Lighthouse</span>
            </NuxtLink>

            <!-- Right Side Icons & Button -->
            <div class="flex items-center gap-2">
                <Profile />
                <ThemePicker />
            </div>
        </div>
    </nav>
</template>

<script setup lang="ts">
import ThemePicker from './ThemePicker.vue';
import Profile from './Profile.vue';

</script>