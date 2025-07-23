<template>
    <div class="w-full shadow px-4 md:px-8">
        <!--- Status Cards --->
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4 p-4 border-border border my-4 rounded-md bg-background">
            <Card class="bg-card text-card-foreground">
                <CardHeader>
                    <CardTitle>
                        Monitors Up
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <p class="text-2xl font-bold">
                        {{ monitorsUp }}
                    </p>
                </CardContent>
            </Card>

            <Card class="bg-card text-card-foreground">
                <CardHeader>
                    <CardTitle>
                        Monitors Down
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <p class="text-2xl font-bold">
                        {{ monitorsDown }}
                    </p>
                </CardContent>
            </Card>
            
            <Card class="bg-card text-card-foreground">
                <CardHeader>
                    <CardTitle>
                        Servers Healthy
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <p class="text-2xl font-bold">
                        {{ serversHealthy }}
                    </p>
                </CardContent>
            </Card>

            <Card class="bg-card text-card-foreground">
                <CardHeader>
                    <CardTitle>
                        Servers Offline
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <p class="text-2xl font-bold">
                        {{ serversOffline }}
                    </p>
                </CardContent>
            </Card>

        </div>
    </div>
</template>

<script setup lang="ts">
definePageMeta({
    middleware: 'auth',
});

import { ref } from 'vue';
import { Card, CardHeader, CardTitle, CardContent } from '~/components/ui/card';

const monitorsUp = ref(0);
const monitorsDown = ref(0);
const serversHealthy = ref(0);
const serversOffline = ref(0);
</script>