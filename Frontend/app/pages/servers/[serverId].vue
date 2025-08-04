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
                            :curve-type="CurveType.Linear" :xDomain="chartTime?.xDomain" />
                    </ClientOnly>
                </CardContent>
            </Card>

            <!--- Memory Usage Chart --->
            <Card>
                <CardHeader>
                    <CardTitle>
                        Memory Usage
                    </CardTitle>
                    <CardDescription>
                        System memory utilization
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <ClientOnly>
                        <AreaChart index="time" :data="memoryChartData" :categories="['memory_cached', 'memory_used']"
                            :show-grid-line="true" :show-legend="true" :show-x-axis="true" :show-y-axis="true"
                            :curve-type="CurveType.Linear" :yDomain="[0, (server?.quickdata?.memory_total || 0) / 1024]"
                            :xDomain="chartTime?.xDomain" />
                    </ClientOnly>
                </CardContent>
            </Card>

            <!--- Disk Usage Chart --->
            <Card>
                <CardHeader>
                    <CardTitle>
                        Disk Usage
                    </CardTitle>
                    <CardDescription>
                        Storage utilization
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <ClientOnly>
                        <AreaChart index="time" :data="storageChartData" :categories="['disk_used']"
                            :show-grid-line="true" :show-legend="true" :show-x-axis="true" :show-y-axis="true"
                            :curve-type="CurveType.Linear" :yDomain="[0, (server?.quickdata?.disk_total || 0) / 1024]"
                            :xDomain="chartTime?.xDomain" />
                    </ClientOnly>
                </CardContent>
            </Card>

            <!--- Disk IO Usage Chart --->
            <Card>
                <CardHeader>
                    <CardTitle>
                        Disk IO Usage
                    </CardTitle>
                    <CardDescription>
                        Storage IO utilization
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <ClientOnly>
                        <AreaChart index="time" :data="storageIOChartData" :categories="['read', 'write']"
                            :show-grid-line="true" :show-legend="true" :show-x-axis="true" :show-y-axis="true"
                            :curve-type="CurveType.Linear"
                            :yDomain="[0, ((server?.quickdata?.disk_io?.read || 0) / 1024) + ((server?.quickdata?.disk_io?.write || 0) / 1024)]"
                            :xDomain="chartTime?.xDomain" />
                    </ClientOnly>
                </CardContent>
            </Card>

            <!--- Network Usage Chart --->
            <Card>
                <CardHeader>
                    <CardTitle>
                        Network Usage
                    </CardTitle>
                    <CardDescription>
                        Network bandwidth utilization
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <ClientOnly>
                        <AreaChart index="time" :data="bandwithChartData" :categories="['download', 'upload']"
                            :show-grid-line="true" :show-legend="true" :show-x-axis="true" :show-y-axis="true"
                            :curve-type="CurveType.Linear" :xDomain="chartTime?.xDomain" />
                    </ClientOnly>
                </CardContent>
            </Card>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import pb from '~/lib/pocketbase';
import { getUTCNowMinus } from '~/lib/pocketbase';
import type { Server, ServerData } from '~/lib/records';
import { ref } from 'vue'
import { CurveType } from '@unovis/ts';

const AreaChart = ref<any>(null);
onMounted(async () => {
    const { AreaChart: LoadedAreaChart } = await import('@/components/ui/chart-area');
    AreaChart.value = LoadedAreaChart;
});

const router = useRouter()
const serverId = router.currentRoute.value.params.serverId as string;
const server = ref<Server>()
const stats = ref<ServerData[]>()
const cpuChartData = ref<Array<{ time: string; cpu_usage: number }>>([])
const bandwithChartData = ref<Array<{ time: string; download: number; upload: number }>>([])
const memoryChartData = ref<Array<{ time: string; memory_cached: number; memory_used: number; memory_total: number }>>([])
const storageChartData = ref<Array<{ time: string; disk_used: number; disk_total: number }>>([])
const storageIOChartData = ref<Array<{ time: string; read: number; write: number }>>([])
const layoutColumns = ref(1) // Layout toggle state
const chartTime = ref<ChartTime>()

const chartTimes: ChartTime[] = [
    {
        xDomain: [0, 60],
        type: "1m",
        hoursAgo: 1
    },
    {
        xDomain: [0, 72],
        type: "5m",
        hoursAgo: 6
    }
]

interface ChartTime {
    xDomain: [number, number];
    type: string;
    hoursAgo: number;
}

const processCpuData = (serverStats: ServerData[]) => {
    return serverStats
        .filter(stat => stat.stats?.cpu_usage !== undefined)
        .map(stat => ({
            time: stat.created.match(/\d{2}:\d{2}:\d{2}/)?.toString() || "",
            cpu_usage: (Math.round((stat.stats?.cpu_usage || 0) * 100) / 100)
        }))
        .reverse()
        .slice(-60) // Show more data points
}

const processBandwithData = (serverStats: ServerData[]) => {
    return serverStats
        .filter(stat => stat.stats?.bandwidth !== undefined)
        .map(stat => ({
            time: stat.created.match(/\d{2}:\d{2}:\d{2}/)?.toString() || "",
            download: stat.stats?.bandwidth?.download || 0,
            upload: stat.stats?.bandwidth?.upload || 0
        }))
        .reverse()
        .slice(-60)
}

const processMemoryData = (serverStats: ServerData[]) => {
    return serverStats
        .filter(stat => stat.stats?.memory_cached !== undefined && stat.stats?.memory_used !== undefined && stat.stats?.memory_total !== undefined)
        .map(stat => ({
            time: stat.created.match(/\d{2}:\d{2}:\d{2}/)?.toString() || "",
            memory_cached: Math.round((stat.stats?.memory_cached || 0) / 1024),
            memory_used: (Math.round(((stat.stats?.memory_used || 0) / 1024) * 100) / 100),
            memory_total: (stat.stats?.memory_total || 0) / 1024
        }))
        .reverse()
        .slice(-60)
}

const processStorageData = (serverStats: ServerData[]) => {
    return serverStats
        .filter(stat => stat.stats?.disk_total !== undefined && stat.stats?.disk_used !== undefined)
        .map(stat => ({
            time: stat.created.match(/\d{2}:\d{2}:\d{2}/)?.toString() || "",
            disk_total: (Math.round(((stat.stats?.disk_total || 0) / 1024) * 100) / 100),
            disk_used: (Math.round(((stat.stats?.disk_used || 0) / 1024) * 100) / 100),
        }))
        .reverse()
        .slice(-60)
}

const processStorageIOData = (serverStats: ServerData[]) => {
    return serverStats
        .filter(stat => stat.stats?.disk_io?.write !== undefined && stat.stats?.disk_io.read !== undefined)
        .map(stat => ({
            time: stat.created.match(/\d{2}:\d{2}:\d{2}/)?.toString() || "",
            write: (Math.round(((stat.stats?.disk_io?.write || 0) / 1024) * 100) / 100),
            read: (Math.round(((stat.stats?.disk_io?.read || 0) / 1024) * 100) / 100),
        }))
        .reverse()
        .slice(-60)
}

onMounted(() => {
    chartTime.value = chartTimes[0]

    pb.collection("servers").getOne(serverId).then((srv) => {
        server.value = srv as unknown as Server
    });

    pb.collection("server_stats").getFullList(60, {
        sort: "-created",
        filter: `server="${serverId}" && type="${chartTime.value?.type}" && created > "${getUTCNowMinus({ hours: chartTime.value?.hoursAgo })}"`
    }).then((datas) => {
        stats.value = datas as unknown as ServerData[]
        cpuChartData.value = processCpuData(stats.value)
        bandwithChartData.value = processBandwithData(stats.value)
        memoryChartData.value = processMemoryData(stats.value)
        storageChartData.value = processStorageData(stats.value)
        storageIOChartData.value = processStorageIOData(stats.value)

        console.log(stats.value[1]?.created)
    })

    setInterval(() => {
        pb.collection("server_stats").getFullList(60, {
            sort: "-created",
            filter: `server="${serverId}" && type="${chartTime.value?.type}" && created > "${getUTCNowMinus({ hours: chartTime.value?.hoursAgo })}"`
        }).then((datas) => {
            stats.value = datas as unknown as ServerData[]
            cpuChartData.value = processCpuData(stats.value)
            bandwithChartData.value = processBandwithData(stats.value)
            memoryChartData.value = processMemoryData(stats.value)

            console.log(stats.value[1]?.created)
        })
    }, 60_000)
})
</script>