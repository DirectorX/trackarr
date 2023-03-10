<template>
    <v-container fluid>
        <v-row align="center">
            <v-col class="pb-0" lg="4" md="4" sm="2" cols="2">
                <h1 class="pt-5 headline font-weight-light">System Logs</h1>
            </v-col>
            <v-col class="text-right pb-0 ml-auto" lg="7" md="7" sm="10" cols="10">
                <v-text-field class="d-inline-flex" v-model="logsSearch" append-icon="mdi-magnify" label="Search"
                    single-line hide-details>
                </v-text-field>
                <v-btn v-on:click="clearLogs($event)" class="ml-3 d-inline-flex">Clear Logs</v-btn>
            </v-col>
        </v-row>
        <v-row>
            <v-col>
                <v-divider class="mb-2"></v-divider>
                <v-data-table :search="logsSearch" calculate-widths :headers="headers" fixed-header
                    :items="filteredMessages()" v-on:update:options="checkSortStatus($event)"
                    :items-per-page.sync="itemsPerPage" class="elevation-1">
                    <template v-slot:item.time="{ item }">
                        <div>
                            {{ item.time | moment("MM/DD/YY, h:mm:ss a")}}
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
                messages: [],
                sorted: false,
                itemsPerPage: 15,
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
                        width: '10%'
                    },
                    {
                        text: 'Level',
                        value: 'level',
                        width: '7%',
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
                        width: '7%',
                        filter: (value) => {
                            if (this.filterComponents.values.length === 0) {
                                return true;
                            }
                            return this.filterComponents.values.includes(value);
                        },
                    },
                    {
                        text: 'Message',
                        width:'76%',
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
                    if (this.logLevels[x].level === level.toUpperCase()) {
                        return this.logLevels[x].color
                    }
                }
                return ""
            },
            clearLogs: function () {
                this.messages = [];
                localStorage.removeItem("logs")

            },
            checkSortStatus: function (event) {
                if (event.sortBy.length === 0) {
                    this.sorted = false;
                } else {
                    this.sorted = true;
                }


            },
            shouldAutoScroll: function () {
                function getScrollPercent() {

                    var h  = document.evaluate("//div[@class='v-data-table__wrapper']", document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue
                    var st = 'scrollTop'
                    var sh = 'scrollHeight'
                    return (h[st]) / ((h[sh]) - h.clientHeight) * 100;
                }


                if (this.sorted || this.itemsPerPage != -1) {
                   
                    return false;
                } else {
                    var percentScrolled = getScrollPercent()
                    if (percentScrolled >= 90) {

                        return true
                    }

                }
                return false;
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

                //Check if autoscroll should occur
                if (this.shouldAutoScroll()) {
                    let table = document.evaluate("//div[@class='v-data-table__wrapper']", document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue
                    table.scrollTop = table.scrollHeight
                }

            }
        }
    };
</script>