<template>
    <div class="w-full shadow px-4 md:px-8">
        <!--- Header --->
        <div class="my-4 flex items-center justify-between">
            <h2 class="text-3xl font-bold tracking-tight">Servers</h2>
            <Form v-slot="{ handleSubmit }" as="" keep-values :validation-schema="formSchema">
                <Dialog>
                    <DialogTrigger as-child>
                        <Button variant="outline">
                            Add Server
                        </Button>
                    </DialogTrigger>
                    <DialogContent class="sm:max-w-[425px]">
                        <DialogHeader>
                            <DialogTitle>Add Server</DialogTitle>
                            <DialogDescription>
                                Input all server details below to create a new server.
                            </DialogDescription>
                        </DialogHeader>

                        <form id="dialogForm" @submit="handleSubmit($event, onSubmit)">
                            <FormField v-slot="{ componentField }" name="name">
                                <FormItem>
                                    <FormLabel>Name</FormLabel>
                                    <FormControl>
                                        <Input type="text" placeholder="Server Name" v-bind="componentField" />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            </FormField>
                            <FormField v-slot="{ componentField }" name="hostname">
                                <FormItem>
                                    <FormLabel>Hostname</FormLabel>
                                    <FormControl>
                                        <Input type="text" placeholder="127.0.0.1" v-bind="componentField" />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            </FormField>
                            <FormField v-slot="{ componentField }" name="port">
                                <FormItem>
                                    <FormLabel>Port</FormLabel>
                                    <FormControl>
                                        <Input type="text" placeholder="8040" v-bind="componentField" />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            </FormField>
                        </form>

                        <DialogFooter>
                            <Button type="submit" form="dialogForm">
                                Save server
                            </Button>
                        </DialogFooter>
                    </DialogContent>
                </Dialog>
            </Form>
        </div>

        <!--- Server List --->
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4 p-4 border-border border my-4 rounded-md bg-background">
            <div v-if="servers.length === 0" class="col-span-1 md:col-span-4">
                <p class="text-muted-foreground text-sm">
                    No servers added yet. Click the button above to add a new server.
                </p>
            </div>
            <!--- System stats for ram, cpu, disk, temp, agent version as a Table --->
            <div class="w-full overflow-x-auto">
                <Table class="min-w-[900px] w-full">
                    <TableHeader>
                        <TableRow>
                            <TableHead>Name</TableHead>
                            <TableHead>Hostname</TableHead>
                            <TableHead>Port</TableHead>
                            <TableHead>Status</TableHead>
                            <TableHead>RAM Usage</TableHead>
                            <TableHead>CPU Usage</TableHead>
                            <TableHead>Disk Usage</TableHead>
                            <TableHead>Temperature</TableHead>
                            <TableHead>Agent Version</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        <template v-for="server in servers" :key="server.id">
                            <TableRow>
                                <TableCell>{{ server.name }}</TableCell>
                                <TableCell>{{ server.hostname }}</TableCell>
                                <TableCell>{{ server.port || 'N/A' }}</TableCell>
                                <TableCell>{{ server.state }}</TableCell>
                                <TableCell>{{ server.quickdata?.ram || 'N/A' }}</TableCell>
                                <TableCell>{{ server.quickdata?.cpu || 'N/A' }}</TableCell>
                                <TableCell>{{ server.quickdata?.disk || 'N/A' }}</TableCell>
                                <TableCell>{{ server.quickdata?.temp || 'N/A' }}</TableCell>
                                <TableCell>{{ server.quickdata?.agent_version || 'N/A' }}</TableCell>
                            </TableRow>
                        </template>
                    </TableBody>
                </Table>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { h } from 'vue'
import * as z from 'zod'
import { toast } from 'vue-sonner'
import pb from '~/lib/pocketbase'
import type { Server } from '~/lib/records'

const formSchema = toTypedSchema(z.object({
    name: z.string().min(2).max(50),
    hostname: z.string().min(2).max(100),
    port: z.string().min(1).max(5).regex(/^\d+$/).optional(),
}))

const servers = ref<Server[]>([])

function onSubmit(values: any) {
    values.users = pb.authStore.record!.id
    values.state = "unknown"
    values.quickdata = {}

    pb.collection('servers').create(values)
        .then((server) => {
            console.log('Server created:', server)
            toast.success('Server created successfully!', {
                description: h('pre', { class: 'mt-2 rounded-md bg-slate-950 p-4' }, h('code', { class: 'text-white' }, JSON.stringify(server, null, 2))),
            })
        })
        .catch((error) => {
            console.error('Error creating server:', error)
            toast.error('Failed to create server', {
                description: h('pre', { class: 'mt-2 rounded-md bg-red-950 p-4' }, h('code', { class: 'text-white' }, error.message + 'ERR: ' + JSON.stringify(values, null, 2))),
            })
        })
}

onMounted(() => {
    pb.collection('servers').getFullList(200, {
        sort: '-created',
    }).then((result) => {
        servers.value = result as unknown as Server[]
        console.log('Servers:', result)
        toast.success('Servers fetched successfully!', {
            description: h('pre', { class: 'mt-2 rounded-md bg-slate-950 p-4' }, h('code', { class: 'text-white' }, JSON.stringify(result, null, 2))),
        })
    }).catch((error) => {
        console.error('Error fetching servers:', error)
        toast.error('Failed to fetch servers', {
            description: h('pre', { class: 'mt-2 rounded-md bg-red-950 p-4' }, h('code', { class: 'text-white' }, error.message)),
        })
    });
});
</script>