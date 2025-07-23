<template>
    <div class="w-full shadow px-4 md:px-8">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4 p-4 border-border border my-4 rounded-md bg-background">
            <Card>
                <CardHeader>
                    <CardTitle class="text-2xl">
                        Main Status Page
                    </CardTitle>
                    <CardDescription>
                        /pages/main
                    </CardDescription>
                </CardHeader>
            </Card>
        </div>
    </div>
</template>

<script setup lang="ts">
definePageMeta({
    middleware: 'auth',
});
</script>