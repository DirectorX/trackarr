<template>
    <v-container fluid>
        <v-row>
            <v-col class="pb-0" lg="5" md="5" sm="5" cols="5">
                <h1 class="pt-5 headline font-weight-light">Pushed Releases</h1>
            </v-col>
            <v-col class="text-right pb-0 ml-auto" lg="6" md="6" sm="7" cols="7">
                <v-text-field class="d-inline-flex" v-model="trackerSearch" append-icon="mdi-magnify" label="Search" single-line hide-details>
                </v-text-field>
            </v-col>
        </v-row>
        <v-row>
            <v-col>
                <v-divider class="mb-2"></v-divider>
                <v-data-table sort-by="age" sort-desc :search="pushedReleasesSearch" disable-sorting
                    :loading="releasesLoading" calculate-widths fixed-header :headers="headersAll" :items="allReleases"
                    :items-per-page="5" class="elevation-1">
                    <template v-if="$vuetify.breakpoint.xs" v-slot:item.release="{ item }">
                        <div class="pl-12">
                            {{ item.release }}
                        </div>
                    </template>
                    <template v-slot:item.age="{ item }">
                        {{ item.age | moment("from", "now")}}
                    </template>
                    <template v-slot:item.pvr="{ item }">
                        {{ item.pvr | capitalize }}
                    </template>
                    <template v-slot:body.append>
                        <tr>
                            <td class="d-none d-sm-table-cell"></td>
                            <td class="d-none d-sm-table-cell"></td>
                            <td
                                :class="{'mt-6 mb-6':$vuetify.breakpoint.xs,'pt-5':$vuetify.breakpoint.smAndUp,'v-data-table__mobile-row':$vuetify.breakpoint.xs,'text-start':!$vuetify.breakpoint.xs}">
                                <v-row>
                                    <v-col class="pt-0 pb-0">
                                        <v-select label="Trackers" prepend-icon="mdi-filter" dense multiple clearable
                                            :items="getTrackers()" v-model="filters.trackersAll">
                                        </v-select>
                                    </v-col>
                                </v-row>

                            </td>
                            <td
                                :class="{'mt-6 mb-6':$vuetify.breakpoint.xs,'pt-5':$vuetify.breakpoint.smAndUp,'v-data-table__mobile-row':$vuetify.breakpoint.xs,'text-start':!$vuetify.breakpoint.xs}">
                                <v-row>
                                    <v-col class="pt-0 pb-0">
                                        <v-select label="PVR" prepend-icon="mdi-filter" dense multiple clearable
                                            :items="getPVR()" v-model="filters.PVRAll">
                                        </v-select>
                                    </v-col>
                                </v-row>
                            </td>
                        </tr>
                    </template>
                </v-data-table>
            </v-col>
        </v-row>
        <v-row>
            <v-col class="pb-0" lg="5" md="5" sm="5" cols="5">
                <h1 class="pt-5 headline font-weight-light">Approved Releases</h1>
            </v-col>
            <v-col class="text-right pb-0 ml-auto" lg="6" md="6" sm="7" cols="7">
                <v-text-field class="d-inline-flex" v-model="trackerSearch" append-icon="mdi-magnify" label="Search" single-line hide-details>
                </v-text-field>
            </v-col>
        </v-row>
        <v-row>
            <v-col>
                <v-divider class="mb-2"></v-divider>
                <v-data-table sort-by="age" sort-desc :search="approvedReleasesSearch" disable-sorting
                    :loading="releasesLoading" calculate-widths fixed-header :headers="headersApproved" :items="approvedReleases"
                    :items-per-page="5" class="elevation-1">
                    <template v-slot:item.age="{ item }">
                        {{ item.age | moment("from", "now")}}
                    </template>
                    <template v-slot:item.pvr="{ item }">
                        {{ item.pvr | capitalize }}
                    </template>
                    <template v-if="$vuetify.breakpoint.xs" v-slot:item.release="{ item }">
                        <div class="pl-12">
                            {{ item.release }}
                        </div>
                    </template>
                    <template v-slot:body.append>
                        <tr>
                            <td class="d-none d-sm-table-cell"></td>
                            <td class="d-none d-sm-table-cell"></td>
                            <td
                                :class="{'mt-6 mb-6':$vuetify.breakpoint.xs,'pt-5':$vuetify.breakpoint.smAndUp,'v-data-table__mobile-row':$vuetify.breakpoint.xs,'text-start':!$vuetify.breakpoint.xs}">
                                <v-row>
                                    <v-col class="pt-0 pb-0">
                                        <v-select label="Trackers" prepend-icon="mdi-filter" dense multiple clearable
                                            :items="getTrackers(approved=true)" v-model="filters.trackersApproved">
                                        </v-select>
                                    </v-col>
                                </v-row>

                            </td>
                            <td
                                :class="{'mt-6 mb-6':$vuetify.breakpoint.xs,'pt-5':$vuetify.breakpoint.smAndUp,'v-data-table__mobile-row':$vuetify.breakpoint.xs,'text-start':!$vuetify.breakpoint.xs}">
                                <v-row>
                                    <v-col class="pt-0 pb-0">
                                        <v-select label="PVR" prepend-icon="mdi-filter" dense multiple clearable
                                            :items="getPVR(approved=true)" v-model="filters.PVRApproved">
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

</template>

<script>
    export default {
        name: 'home',
        data() {
            return {
                allReleases: [],
                approvedReleases: [],
                pushedReleasesSearch: '',
                releasesLoading: true,
                approvedReleasesSearch: '',
                filters: {
                    trackersAll: [],
                    PVRAll: [],
                    trackersApproved: [],
                    PVRApproved: []
                }
            }
        },
        methods: {
            fetchReleases: function () {
                this.$axios.get('/releases', {
                    params: {
                        apikey: this.CORE_API_KEY
                    }
                }).then(
                    response => {
                        for (let i = 0; i < response.data.length; i++) {
                            this.allReleases.push({
                                age: response.data[i].CreatedAt,
                                release: response.data[i].Name,
                                pvr: response.data[i].PvrName,
                                tracker: response.data[i].TrackerName
                            });

                            if (response.data[i].Approved) {
                                this.approvedReleases.push({
                                    age: response.data[i].CreatedAt,
                                    release: response.data[i].Name,
                                    pvr: response.data[i].PvrName,
                                    tracker: response.data[i].TrackerName
                                })
                            }
                        }
                        this.releasesLoading = false;
                    })
            },
            getTrackers: function (approved = false) {
                if (!approved) {
                    return [...new Set(this.allReleases.map(item => item.tracker))]
                }
                return [...new Set(this.approvedReleases.map(item => item.tracker))]
            },
            getPVR: function (approved = false) {
                if (!approved) {
                    return [...new Set(this.allReleases.map(item => item.pvr))]
                }
                return [...new Set(this.approvedReleases.map(item => item.pvr))]

            }
        },
        beforeDestroy: function () {
            // unsubscribe from releases topic
            if (this.$socket.readyState === WebSocket.OPEN) {
                this.$socket.sendObj({
                    type: 'unsubscribe',
                    data: 'releases'
                });
            }
        },
        mounted: function () {
            // retrieve releases
            this.fetchReleases();

            // subscribe to releases topic when socket is already open
            if (this.$socket.readyState === WebSocket.OPEN) {
                this.$socket.sendObj({
                    type: 'subscribe',
                    data: 'releases'
                });
            } else {
                // subscribe to releases topic when socket is connected
                this.$options.sockets.onopen = () => {
                    this.$socket.sendObj({
                        type: 'subscribe',
                        data: 'releases'
                    });
                };
            }

            // set message handler
            this.$options.sockets.onmessage = (message) => {
                // parse message
                let event = JSON.parse(message.data);

                // ignore irrelevant messages
                if (!event.type || event.type !== 'release')
                    return;

                // handle release message
                this.allReleases.push({
                    age: event.data.CreatedAt,
                    release: event.data.Name,
                    pvr: event.data.PvrName,
                    tracker: event.data.TrackerName
                });

                if (event.data.Approved) {
                    this.approvedReleases.push({
                        age: event.data.CreatedAt,
                        release: event.data.Name,
                        pvr: event.data.PvrName,
                        tracker: event.data.TrackerName
                    });
                }
            }
        },
        computed: {
            headersAll() {
                return [{
                        text: 'Age',
                        value: 'age',
                    },
                    {
                        text: 'Release',
                        value: 'release'
                    },
                    {
                        text: 'Tracker',
                        value: 'tracker',
                        filter: (value) => {
                            if (this.filters.trackersAll.length === 0) {
                                return true;
                            }
                            return this.filters.trackersAll.includes(value);
                        },
                    },
                    {
                        text: 'PVR',
                        value: 'pvr',
                        filter: (value) => {
                            if (this.filters.PVRAll.length === 0) {
                                return true;
                            }
                            return this.filters.PVRAll.includes(value.toLowerCase());
                        },
                    }
                ]
            },
            headersApproved() {
                return [{
                        text: 'Age',
                        value: 'age',
                    },
                    {
                        text: 'Release',
                        value: 'release'
                    },
                    {
                        text: 'Tracker',
                        value: 'tracker',
                        filter: (value) => {
                            if (this.filters.trackersApproved.length === 0) {
                                return true;
                            }
                            return this.filters.trackersApproved.includes(value);
                        },
                    },
                    {
                        text: 'PVR',
                        value: 'pvr',
                        filter: (value) => {
                            if (this.filters.PVRApproved.length === 0) {
                                return true;
                            }
                            return this.filters.PVRApproved.includes(value.toLowerCase());
                        },
                    }
                ]
            }


        }

    };
</script>