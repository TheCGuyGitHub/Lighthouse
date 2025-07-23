<template>
    <DropdownMenu>
        <DropdownMenuTrigger as-child class="flex items-center gap-2">
            <Button variant="ghost">
                <User class="h-[1.2rem] w-[1.2rem]" />
            </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
            <DropdownMenuLabel>
                {{ username }}
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem>
                Settings
            </DropdownMenuItem>
            <DropdownMenuItem @click="pb.authStore.clear()">
                Logout
            </DropdownMenuItem>
        </DropdownMenuContent>
    </DropdownMenu>
</template>

<script setup lang="ts">
import { User } from 'lucide-vue-next';
import { ref } from 'vue';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuTrigger } from '@/components/ui/dropdown-menu';
import { Button } from '@/components/ui/button';
import pb from '@/lib/pocketbase';

const username = ref('');

onMounted(() => {
    if (pb.authStore.isValid) {
        username.value = pb.authStore.record?.name || pb.authStore.record?.email;
    } else {
        username.value = 'Guest';
    }
});
</script>