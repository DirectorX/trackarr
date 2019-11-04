<template>
  <v-container fluid>
    <v-row>
      <v-col class="pb-0" cols="9">
        <h1 class="pt-5 headline font-weight-light">Pushed Releases</h1>
      </v-col>
      <v-col class="pb-0" cols="3">
        <v-text-field v-model="pushedReleasesSearch" append-icon="mdi-magnify" label="Search" single-line hide-details></v-text-field>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-divider class="mb-2"></v-divider>
        <v-data-table :search="pushedReleasesSearch" disable-sorting :loading="releasesLoading" calculate-widths :headers="headers"
          :items="allReleases" :items-per-page="5" class="elevation-1">
          <template v-slot:item.age="{ item }">
            {{ item.age | moment("from", "now")}}
          </template>
          <template v-slot:item.pvr="{ item }">
            {{ item.pvr | capitalize }}
          </template>
        </v-data-table>
      </v-col>
    </v-row>
    <v-row>
      <v-col class="pb-0" cols="9">
        <h1 class="pt-5 headline font-weight-light">Approved Releases</h1>
      </v-col>
       <v-col class="pb-0" cols="3">
        <v-text-field v-model="approvedReleasesSearch" append-icon="mdi-magnify" label="Search" single-line hide-details></v-text-field>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-divider class="mb-2"></v-divider>
        <v-data-table :search="approvedReleasesSearch" disable-sorting :loading="releasesLoading" calculate-widths :headers="headers"
          :items="approvedReleases" :items-per-page="5" class="elevation-1">
          <template v-slot:item.age="{ item }">
            {{ item.age | moment("from", "now")}}
          </template>
          <template v-slot:item.pvr="{ item }">
            {{ item.pvr | capitalize }}
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
        headers: [{
            text: 'Age',
            value: 'age',
          },
          {
            text: 'Release',
            value: 'release'
          },
          {
            text: 'Tracker',
            value: 'tracker'
          },
          {
            text: 'PVR',
            value: 'pvr'
          }
        ],
        allReleases: [],
        approvedReleases: [],
        pushedReleasesSearch: '',
        releasesLoading: true,
        approvedReleasesSearch: ''
      }
    },
    methods: {
      fetchReleases: function () {
        this.$axios.get(process.env.VUE_APP_RELEASE_URL).then(
          response => {
            for (let i = 0; i < response.data.length; i++) {
              this.allReleases.push({
                age: response.data[i].CreatedAt,
                release: response.data[i].Name,
                pvr: response.data[i].PvrName,
                tracker: response.data[i].TrackerName
              })

              if(response.data[i].Approved){
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
      }
    },
    mounted: function () {
      this.fetchReleases()
    },
    filters: {
      capitalize: function (value) {
        if (!value) return ''
        value = value.toString()
        return value.charAt(0).toUpperCase() + value.slice(1)
      }
    }

  };
</script>