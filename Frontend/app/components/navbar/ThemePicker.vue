<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { Sun, Moon, Laptop } from 'lucide-vue-next'

const colorMode = useColorMode()

const setTheme = (theme: 'light' | 'dark' | 'system') => {
    colorMode.preference = theme
}

</script>

<template>
    <DropdownMenu>
        <DropdownMenuTrigger as-child>
            <Button variant="ghost" size="icon">
                <Sun v-if="colorMode.value === 'light'" class="h-[1.2rem] w-[1.2rem]" />
                <Moon v-else-if="colorMode.value === 'dark'" class="h-[1.2rem] w-[1.2rem]" />
                <Laptop v-else class="h-[1.2rem] w-[1.2rem]" />
                <span class="sr-only">Toggle theme</span>
            </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
            <DropdownMenuItem @click="setTheme('light')">
                <Sun class="mr-2 h-4 w-4" />
                <span>Light</span>
            </DropdownMenuItem>
            <DropdownMenuItem @click="setTheme('dark')">
                <Moon class="mr-2 h-4 w-4" />
                <span>Dark</span>
            </DropdownMenuItem>
            <DropdownMenuItem @click="setTheme('system')">
                <Laptop class="mr-2 h-4 w-4" />
                <span>System</span>
            </DropdownMenuItem>
        </DropdownMenuContent>
    </DropdownMenu>
</template>