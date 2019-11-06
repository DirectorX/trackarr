<template>
  <v-container fluid>
    <v-row>
      <v-col class="pb-0" cols="9">
        <h1 class="pt-5 headline font-weight-light">System Logs</h1>
      </v-col>
      <v-col class="pb-0" cols="3">
        <v-text-field v-model="logsSearch" append-icon="mdi-magnify" label="Search" single-line hide-details></v-text-field>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-divider class="mb-2"></v-divider>
        <v-data-table :search="logsSearch" disable-sorting calculate-widths :headers="headers"
          :items="messages" :items-per-page="5" class="elevation-1">
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
        headers: [{
            text: 'Timestamp',
            value: 'timestamp',
          },
          {
            text: 'Level',
            value: 'level'
          },
          {
            text: 'Component',
            value: 'component'
          },
          {
            text: 'Message',
            value: 'message'
          }
        ],
        messages: [],
        logsSearch: '',

      }
    },
    methods: {
    },
    mounted: function () {
      this.$options.sockets.onmessage = (data) =>{
        console.log(data)
        this.messages.push({
          timestamp: "foo",
          level:"WARN",
          component:"Test Component",
          message:"Test Message"
        })
      }
    }

  };
</script>
