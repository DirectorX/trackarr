<template>
    <div>
        <v-container fluid>
            <v-row align="center">
                <v-col class="pb-0" lg="4" md="4" sm="4" cols="4">
                    <h1 class="pt-5 headline font-weight-light">Tracker Status</h1>
                </v-col>
                <v-col class="text-right pb-0 ml-auto" lg="7" md="7" sm="8" cols="8">
                    <v-text-field class="d-inline-flex" v-model="trackerSearch" append-icon="mdi-magnify" label="Search"
                                  single-line hide-details>
                    </v-text-field>
                </v-col>
            </v-row>
            <v-row>
                <v-col>
                    <v-divider class="mb-2"></v-divider>
                    <v-data-table :search="trackerSearch" disable-sorting calculate-widths :headers="headers"
                                  :items="filteredTrackers()" :items-per-page="15 " class="elevation-1">

                        <template v-slot:item.status="{ item }">
                            <div>
                                <status-indicator v-if="item.status == 'online'" class="mr-3" status="positive" pulse/>
                                <status-indicator v-if="item.status == 'offline'" class="mr-3" status="negative"
                                                  pulse/>
                                {{ item.status | capitalize }}
                            </div>

                        </template>
                        <template v-slot:item.last_announced="{ item }">
                            <div>
                                <span v-if="item.last_announced !== null">
                                    {{ item.last_announced | moment("from", "now") | capitalize }}
                                </span>
                                <span v-else>
                                    Never
                                </span>
                            </div>
                        </template>
                        <template v-slot:body.append>
                            <tr>
                                <td class="d-none d-sm-table-cell"></td>
                                <td
                                        :class="{'mt-6 mb-6':$vuetify.breakpoint.xs,'pt-5':$vuetify.breakpoint.smAndUp,'v-data-table__mobile-row':$vuetify.breakpoint.xs,'text-start':!$vuetify.breakpoint.xs}">
                                    <v-row>
                                        <v-col class="pt-0 pb-0">
                                            <v-select label="Tracker Status" prepend-icon="mdi-filter" dense clearable
                                                      :items="trackerStatus" v-model="filterTrackerStatus">
                                            </v-select>
                                        </v-col>
                                    </v-row>

                                </td>
                            </tr>
                        </template>
                    </v-data-table>
                </v-col>
            </v-row>
        </v-container>
        <v-footer absolute dark padless>
            <v-col class="text-right" cols="12">
                <v-btn v-on:click="checkForUpdate()" color="secondary" class="ma-2">Check For Update</v-btn>
                <label><strong>Trackarr {{currentVersion}}</strong></label>

            </v-col>


        </v-footer>
    </div>

</template>


<script>
    import {StatusIndicator} from 'vue-status-indicator';

    export default {
        name: 'status',
        components: {
            StatusIndicator,
        },
        data() {
            return {
                trackers: [],
                trackerSearch: '',
                filterTrackerStatus: '',
                trackerStatus: [{
                    text: "Online",
                    value: "online",
                },
                    {
                        text: "Offline",
                        value: "offline",
                    }

                ],
                currentVersion: this.CORE_APP_VERSION
            }
        },
        computed: {
            headers() {
                return [{
                    text: 'Tracker',
                    value: 'tracker',
                },
                    {
                        text: 'Status',
                        value: 'status',
                        filter: (value) => {
                            if (!this.filterTrackerStatus) {
                                return true;
                            }
                            return this.filterTrackerStatus == value;
                        },
                    },
                    {
                        text: 'Last Announced',
                        value: 'last_announced'
                    }
                ]
            }
        },
        methods: {
            filteredTrackers: function () {
                return this.trackers.filter(item => {
                    if (this.filterTrackerStatus) {
                        if (item.status == this.filterTrackerStatus) {
                            return true
                        }
                        return false

                    }
                    return true;
                })

            },
            fetchTrackerStatuses: function () {
                this.$axios.get('/irc/status', {
                    params: {
                        apikey: this.CORE_API_KEY
                    }
                }).then(
                    response => {
                        for (const [key, value] of Object.entries(response.data)) {
                            let index = this.trackers.map(item => item.tracker).indexOf(key);
                            if (index === -1) {
                                this.trackers.push({
                                    tracker: key,
                                    status: value.connected === true ? 'online' : 'offline',
                                    last_announced: value.last_announced !== '' ? value.last_announced : 'Never'
                                })
                            } else {
                                this.trackers[index].status = value.connected === true ? 'online' : 'offline';
                                this.trackers[index].last_announced = value.last_announced !== '' ? value.last_announced : null
                            }
                        }
                    })
            },
            checkForUpdate: function () {
                this.$axios.get('/update/status', {
                    params: {
                        apikey: this.CORE_API_KEY
                    }
                }).then(response => {
                    if (response.data.update_available) {
                        this.$toastr.Add({
                            msg: "Newer version available: " + response.data.latest_version,
                            position: 'toast-top-right'
                        })
                    } else {
                        this.$toastr.Add({
                            msg: "You're already on the latest version!",
                            position: 'toast-top-right'
                        })
                    }
                })
            }
        },
        beforeDestroy: function () {
            clearInterval(this.interval)
        },
        mounted: function () {
            // retrieve releases
            this.fetchTrackerStatuses();
            this.interval = setInterval(() => {
                this.fetchTrackerStatuses()
            }, 60000);
        }
    };
</script>