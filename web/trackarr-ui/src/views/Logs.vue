<template>
  <v-container fluid>
    <v-row>
      <v-col class="pb-0" cols="4">
        <h1 class="pt-5 headline font-weight-light">System Logs</h1>
      </v-col>
      <v-col class="pb-0" offset-lg="5" offset-md="5" offset-sm="5" offset-xs="5" cols="3">
        <v-text-field v-model="logsSearch" append-icon="mdi-magnify" label="Search" single-line hide-details></v-text-field>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-divider class="mb-2"></v-divider>
        <v-data-table :search="logsSearch" disable-sorting calculate-widths :headers="headers"
          :items="messages" :items-per-page="15 " class="elevation-1">
          <template v-slot:item.level="{ item }">
            {{ item.level.toUpperCase() }}
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
        headers: [{
            text: 'Timestamp',
            value: 'time',
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
      if(localStorage.logs){ 
          this.messages = JSON.parse(localStorage.logs)
      }

      this.$options.sockets.onopen = () => {
        this.$socket.sendObj({type: 'subscribe', data: 'logs'})
      }
      
      this.$options.sockets.onmessage = (message) => {
        this.messages.push(JSON.parse(message.data).data)
        if(localStorage.logs){
          let data = JSON.parse(localStorage.logs)
          if(data.length > 500){
            localStorage.logs = null
          }
        }
        localStorage.logs = JSON.stringify(this.messages)
      }
    }

  };
</script>
