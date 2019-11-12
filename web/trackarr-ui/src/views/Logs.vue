<template>
    <v-container fluid>
        <v-row align="center">
            <v-col class="pb-0" lg="4" md="4" sm="4" cols="4">
                <h1 class="pt-5 headline font-weight-light">System Logs</h1>
            </v-col>
            <v-col class="text-right pb-0 ml-auto" lg="7" md="7" sm="8" cols="8">
                <v-text-field class="d-inline-flex" v-model="logsSearch" append-icon="mdi-magnify" label="Search" single-line hide-details>
                </v-text-field>
                <v-btn v-on:click="clearLogs()" class="ml-3 d-inline-flex">Clear Logs</v-btn>
            </v-col>
        </v-row>
        <v-row>
            <v-col>
                <v-divider class="mb-2"></v-divider>
                <v-data-table :search="logsSearch" disable-sorting calculate-widths :headers="headers"
                    :items="filteredMessages()" :items-per-page="15 " class="elevation-1">
                    <template v-slot:item.time="{ item }">
                        <div>
                            {{ item.time }}
                        </div>
                    </template>
                    <template v-slot:item.level="{ item }">
                        <div :style="{color:getLogColor(item.level)}">
                            {{ item.level.toUpperCase() }}
                        </div>
                    </template>
                    <template v-if="$vuetify.breakpoint.xs" v-slot:item.message="{ item }">
                        <div class="pl-12">
                            {{ item.message }}
                        </div>
                    </template>
                    <template v-slot:item.component="{ item }">
                        {{ item.component }}
                    </template>
                    <template v-slot:body.append>
                        <tr>
                            <td class="d-none d-sm-table-cell"></td>
                            <td
                                :class="{'mt-6 mb-6':$vuetify.breakpoint.xs,'pt-5':$vuetify.breakpoint.smAndUp,'v-data-table__mobile-row':$vuetify.breakpoint.xs,'text-start':!$vuetify.breakpoint.xs}">
                                <v-row>
                                    <v-col class="pt-0 pb-0">
                                        <v-select label="Log Level" prepend-icon="mdi-filter" dense multiple clearable
                                            :items="logLevels" item-text="level" item-value="level"
                                            v-model="filterLevels.values">
                                        </v-select>
                                    </v-col>
                                </v-row>

                            </td>
                            <td
                                :class="{'mt-6 mb-6':$vuetify.breakpoint.xs,'pt-5':$vuetify.breakpoint.smAndUp,'v-data-table__mobile-row':$vuetify.breakpoint.xs,'text-start':!$vuetify.breakpoint.xs}">
                                <v-row>
                                    <v-col class="pt-0 pb-0">
                                        <v-select label="Component" prepend-icon="mdi-filter" dense multiple clearable
                                            :items="getComponents()" v-model="filterComponents.values">
                                        </v-select>
                                    </v-col>
                                </v-row>
                            </td>
                            <td class="d-none d-sm-table-cell"></td>
                        </tr>
                    </template>
                </v-data-table>
            </v-col>
        </v-row>
    </v-container>
</template>


<script>
    export default {
        name: 'logs',
        data() {
            return {
                messages: [{
                    time:"foo",
                    level:"INFO",
                    component:"app",
                    message:"This is a test message"
                }],
                logsSearch: '',
                filterLevels: {
                    values: []
                },
                filterComponents: {
                    values: []
                },
                logLevels: [{
                        level: "TRACE",
                        color: "#808080",
                    },
                    {
                        level: "DEBUG",
                        color: "#00AAAA",
                    },
                    {
                        level: "INFO",
                        color: "#00AA00",
                    },
                    {
                        level: "WARN",
                        color: "#AAAA00",
                    },
                    {
                        level: "ERROR",
                        color: "#AA0000",
                    }
                ]

            }
        },
        computed: {
            headers() {
                return [{
                        text: 'Timestamp',
                        value: 'time',
                    },
                    {
                        text: 'Level',
                        value: 'level',
                        filter: (value) => {
                            if (this.filterLevels.values.length === 0) {
                                return true;
                            }
                            return this.filterLevels.values.includes(value.toUpperCase());
                        },
                    },
                    {
                        text: 'Component',
                        value: 'component',
                        filter: (value) => {
                            if (this.filterComponents.values.length === 0) {
                                return true;
                            }
                            return this.filterComponents.values.includes(value);
                        },
                    },
                    {
                        text: 'Message',
                        value: 'message'
                    }
                ]
            }
        },
        methods: {
            getComponents: function () {
                return [...new Set(this.messages.map(item => item.component))]
            },
            filteredMessages: function () {
                return this.messages.filter(item => {
                    if (this.filterComponents.length > 0) {
                        if (!this.filterComponents.includes(item.component)) {
                            return false;
                        }
                        if (this.filterLevels.length > 0 && !this.filterLevels.includes(item.level)) {
                            return false;
                        }
                    } else {
                        if (this.filterLevels.length > 0 && !this.filterLevels.includes(item.level)) {
                            return false;
                        }
                    }
                    return true;
                })

            },
            getLogColor: function (level) {
                for (let x = 0; x < this.logLevels.length; x++) {
                    if (this.logLevels[x].level === level) {
                        return this.logLevels[x].color
                    }
                }
                return ""
            },
            clearLogs: function(){
                this.messages = []
                localStorage.removeItem("logs")

            }
        },
        beforeDestroy: function () {
            // unsubscribe from logs topic
            if (this.$socket.readyState === WebSocket.OPEN) {
                this.$socket.sendObj({
                    type: 'unsubscribe',
                    data: 'logs'
                });
            }
        },
        mounted: function () {
            // load messages from store
            if (localStorage.logs) {
                this.messages = JSON.parse(localStorage.logs);
            }

            // subscribe to logs topic when socket is already open
            if (this.$socket.readyState === WebSocket.OPEN) {
                this.$socket.sendObj({
                    type: 'subscribe',
                    data: 'logs'
                });
            } else {
                // subscribe to logs topic when socket is connected
                this.$options.sockets.onopen = () => {
                    this.$socket.sendObj({
                        type: 'subscribe',
                        data: 'logs'
                    });
                };
            }

            // set message handler
            this.$options.sockets.onmessage = (message) => {
                // parse message
                let event = JSON.parse(message.data);

                // ignore irrelevant messages
                if (!event.type || event.type !== 'log')
                    return;

                // handle log message
                this.messages.push(event.data);
                if (localStorage.logs) {
                    let data = JSON.parse(localStorage.logs);
                    if (data.length > 500) {
                        localStorage.logs = null;
                    }
                }

                localStorage.logs = JSON.stringify(this.messages);
            }
        }
    };
</script>