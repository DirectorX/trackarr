<template>
  <v-container fluid>
    <v-row>
      <v-col class="pb-0" cols="4">
        <h1 class="pt-5 headline font-weight-light">System Logs</h1>
      </v-col>
      <v-col class="pb-0" offset-lg="5" offset-md="5" offset-sm="5" offset-xs="5" cols="3">
        <v-text-field v-model="logsSearch" append-icon="mdi-magnify" label="Search" single-line hide-details>
        </v-text-field>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-divider class="mb-2"></v-divider>
        <v-data-table :search="logsSearch" disable-sorting calculate-widths :headers="headers" :items="filteredMessages()"
          :items-per-page="15 " class="elevation-1">
          <template v-slot:item.level="{ item }">
            {{ item.level.toUpperCase() }}
          </template>
          <template v-slot:body.append>
            <tr>
              <td class="d-none d-sm-flex"></td>
              <td :class="{'mt-6 mb-6':$vuetify.breakpoint.xs,'pt-5':$vuetify.breakpoint.smAndUp,'v-data-table__mobile-row':$vuetify.breakpoint.xs,'text-start':!$vuetify.breakpoint.xs}">
                <v-row> 
                  <v-col class="pt-0 pb-0">
                      <v-select label="Log Level" prepend-icon="mdi-filter" dense multiple clearable
                  :items="logLevels"
                  
                  :value="filterLevels">
                  </v-select>
                  </v-col>
                </v-row>
              
              </td>
              <td :class="{'mt-6 mb-6':$vuetify.breakpoint.xs,'pt-5':$vuetify.breakpoint.smAndUp,'v-data-table__mobile-row':$vuetify.breakpoint.xs,'text-start':!$vuetify.breakpoint.xs}">
                <v-row>
                  <v-col class="pt-0 pb-0">
                      <v-select label="Component" prepend-icon="mdi-filter" dense multiple clearable
                  :items="getComponents()"
                  :value="filterComponents">
                  </v-select>
                  </v-col>
                </v-row>
              </td>
              <td class="d-none d-sm-flex"></td>
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
        filterLevels: [],
        filterComponents: [],
        logLevels: ["TRACE","DEBUG","INFO","WARN","ERROR","FATAL"]

      }
    },
    methods: {
      getComponents: function(){
          return [...new Set(this.messages.map(item => item.component))]
      },
      filteredMessages: function(){
          return this.messages.filter(item => {
              if(this.filterComponents.length > 0){
                if (!this.filterComponents.includes(item.component)){
                  return false
                }
                if (this.filterLevels.length > 0 && !this.filterLevels.includes(item.level)){
                  return false;
                }
              }
              else{
                if (this.filterLevels.length > 0 && !this.filterLevels.includes(item.level)){
                  return false;
                }
              }
              return true;
          })

      }
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