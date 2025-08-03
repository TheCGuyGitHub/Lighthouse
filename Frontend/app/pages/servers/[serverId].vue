<template>
    <div class="w-full shadow px-4 md:px-8">
        <!--- Header --->
        <div class="my-4 flex items-center justify-between">
            <div
                class="flex items-center justify-between w-full p-4 border-border border my-4 rounded-md bg-background">
                <h2 class="text-3xl font-bold tracking-tight">{{ server?.name }}</h2>

                <!-- Layout Toggle -->
                <div class="flex items-center gap-2">
                    <span class="text-sm text-gray-400">Layout:</span>
                    <button @click="layoutColumns = 1" :class="[
                        'px-3 py-1 text-xs rounded transition-colors',
                        layoutColumns === 1
                            ? 'bg-blue-600 text-white'
                            : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                    ]">
                        1 Column
                    </button>
                    <button @click="layoutColumns = 2" :class="[
                        'px-3 py-1 text-xs rounded transition-colors',
                        layoutColumns === 2
                            ? 'bg-blue-600 text-white'
                            : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                    ]">
                        2 Columns
                    </button>
                </div>
            </div>
        </div>

        <!--- Charts Container --->
        <div :class="[
            'grid gap-6',
            layoutColumns === 1 ? 'grid-cols-1' : 'grid-cols-1 lg:grid-cols-2'
        ]">
            <!--- CPU Usage Chart --->
            <Card>
                <CardHeader>
                    <CardTitle>
                        CPU Usage
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <ClientOnly>
                        <AreaChart index="time" :data="cpuChartData" :categories="['cpu_usage']" :colors="['cyan']"
                            :show-grid-line="true" :show-legend="true" :show-x-axis="true" :show-y-axis="true"
                            :curve-type="CurveType.Linear" />
                    </ClientOnly>
                </CardContent>
            </Card>

            <!--- Memory Usage Chart (Placeholder) --->
            <div class="p-6 border-border border rounded-md bg-background">
                <div class="mb-6">
                    <h3 class="text-xl font-semibold text-white mb-1">Memory Usage</h3>
                    <p class="text-sm text-gray-400">System memory utilization</p>
                </div>
                <div class="h-64 flex items-center justify-center text-gray-400 bg-gray-900 rounded-lg">
                    Memory chart coming soon...
                </div>
            </div>

            <!--- Disk Usage Chart (Placeholder) --->
            <div class="p-6 border-border border rounded-md bg-background">
                <div class="mb-6">
                    <h3 class="text-xl font-semibold text-white mb-1">Disk Usage</h3>
                    <p class="text-sm text-gray-400">Storage utilization</p>
                </div>
                <div class="h-64 flex items-center justify-center text-gray-400 bg-gray-900 rounded-lg">
                    Disk chart coming soon...
                </div>
            </div>

            <!--- Network Usage Chart (Placeholder) --->
            <div class="p-6 border-border border rounded-md bg-background">
                <div class="mb-6">
                    <h3 class="text-xl font-semibold text-white mb-1">Network Usage</h3>
                    <p class="text-sm text-gray-400">Network bandwidth utilization</p>
                </div>
                <div class="h-64 flex items-center justify-center text-gray-400 bg-gray-900 rounded-lg">
                    Network chart coming soon...
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import pb from '~/lib/pocketbase';
import type { Server, ServerData } from '~/lib/records';
import { ref } from 'vue'
import { CurveType } from '@unovis/ts';
// import { AreaChart } from '@/components/ui/chart-area';

const AreaChart = ref<any>(null);

onMounted(async () => {
    // The dynamic import is inside onMounted, so it only runs on the client.
    const { AreaChart: LoadedAreaChart } = await import('@/components/ui/chart-area');
    AreaChart.value = LoadedAreaChart;
});

const router = useRouter()
const serverId = router.currentRoute.value.params.serverId as string;
const server = ref<Server>()
const stats = ref<ServerData[]>()
const cpuChartData = ref<Array<{ time: string; cpu_usage: number }>>([])
const layoutColumns = ref(1) // Layout toggle state

// Data processing
const processCpuData = (serverStats: ServerData[]) => {
    return serverStats
        .filter(stat => stat.stats?.cpu_usage !== undefined)
        .map(stat => ({
            time: stat.created,
            cpu_usage: (Math.round((stat.stats?.cpu_usage || 0) * 100) / 100)
        }))
        .reverse()
        .slice(-60) // Show more data points
}

onMounted(() => {
    pb.collection("servers").getOne(serverId).then((srv) => {
        server.value = srv as unknown as Server
    });

    pb.collection("server_stats").getFullList(100, {
        sort: "-created",
        filter: `server="${serverId}" && type="1m"`
    }).then((datas) => {
        stats.value = datas as unknown as ServerData[]
        cpuChartData.value = processCpuData(stats.value)
    })
})
</script>