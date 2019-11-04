<template>
  <v-container fluid>
    <v-row>
        <v-col>
        <h2 class="headline font-weight-light">All Releases</h2>
        <v-divider class="mb-2"></v-divider>
        <v-data-table :headers="headers" :items="allReleases" :items-per-page="5" class="elevation-1"></v-data-table>
      </v-col>
       <v-col>
        <h2 class="headline font-weight-light">Approved Releases</h2>
        <v-divider class="mb-2"></v-divider>
        <v-data-table :headers="headers" :items="allReleases" :items-per-page="5" class="elevation-1"></v-data-table>
      </v-col>
    </v-row>
     <!-- <v-row>
    <h2>Tracker Statuses</h2>
    </v-row>
    <v-row>
      <v-col>
        <v-card
    class="mx-auto"
    max-width="344"
    outlined
  >
    <v-list-item three-line>
      <v-list-item-content>
        <div class="overline mb-4">OVERLINE</div>
        <v-list-item-title class="headline mb-1">BeyondHD</v-list-item-title>
        <v-list-item-subtitle>Greyhound divisely hello coldly fonwderfully</v-list-item-subtitle>
      </v-list-item-content>

      <v-list-item-avatar
        tile
        size="80"
        color="grey"
      ></v-list-item-avatar>
    </v-list-item>
  </v-card>
      </v-col>
        <v-col>
        <v-card
    class="mx-auto"
    max-width="344"
    outlined
  >
    <v-list-item three-line>
      <v-list-item-content>
        <div class="overline mb-4">OVERLINE</div>
        <v-list-item-title class="headline mb-1">BTN</v-list-item-title>
        <v-list-item-subtitle>Greyhound divisely hello coldly fonwderfully</v-list-item-subtitle>
      </v-list-item-content>

      <v-list-item-avatar
        tile
        size="80"
        color="grey"
      ></v-list-item-avatar>
    </v-list-item>
  </v-card>
      </v-col>
        <v-col>
        <v-card
    class="mx-auto"
    max-width="344"
    outlined
  >
    <v-list-item three-line>
      <v-list-item-content>
        <div class="overline mb-4">OVERLINE</div>
        <v-list-item-title class="headline mb-1">HD-Torrents</v-list-item-title>
        <v-list-item-subtitle>Greyhound divisely hello coldly fonwderfully</v-list-item-subtitle>
      </v-list-item-content>

      <v-list-item-avatar
        tile
        size="80"
        color="grey"
      ></v-list-item-avatar>
    </v-list-item>
  </v-card>
      </v-col>
        <v-col>
        <v-card
    class="mx-auto"
    max-width="344"
    outlined
  >
    <v-list-item three-line>
      <v-list-item-content>
        <div class="overline mb-4">OVERLINE</div>
        <v-list-item-title class="headline mb-1">IPTorrents</v-list-item-title>
        <v-list-item-subtitle>Greyhound divisely hello coldly fonwderfully</v-list-item-subtitle>
      </v-list-item-content>

      <v-list-item-avatar
        tile
        size="80"
        color="grey"
      ></v-list-item-avatar>
    </v-list-item>
  </v-card>
      </v-col>
        <v-col>
        <v-card
    class="mx-auto"
    max-width="344"
    outlined
  >
    <v-list-item three-line>
      <v-list-item-content>
        <div class="overline mb-4">OVERLINE</div>
        <v-list-item-title class="headline mb-1">PassThePopcorn</v-list-item-title>
        <v-list-item-subtitle>Greyhound divisely hello coldly fonwderfully</v-list-item-subtitle>
      </v-list-item-content>

      <v-list-item-avatar
        tile
        size="80"
        color="grey"
      ></v-list-item-avatar>
    </v-list-item>
  </v-card>
      </v-col>
    </v-row> -->
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
          }
        ],
        allReleases: [
            {
                "age":"17h",
                "release": "Fuses.mkv",
                "tracker": "PtP",
                "PvrName": "radarr",
                "Approved": false
            },
            {
                
                "release": "Mondscheinkinder [2006] - AQ.mkv",
                "tracker": "PtP",
                "age":"13h"
            },
            {
                
                "release": "Mondscheinkinder [2006] - AQ.mkv",
                "tracker": "PtP",
                "age":"13h"
            },
            {
                
                "release": "Mondscheinkinder [2006] - AQ.mkv",
                "tracker": "PtP",
                "age":"13h"
            },
            {
                
                "release": "Mondscheinkinder [2006] - AQ.mkv",
                "tracker": "PtP",
                "age":"13h"
            },
            ],
      }
    },
    methods: {
      fetchReleases: function(){
        this.$axios.get("http://trackarr.cloudbox.media:7337/api/releases?apikey=RGqVz47Aioa5JxwmuwAEJY").then(response =>{
                    console.log(response);
                })
      }
    },
    mounted: function(){
      this.fetchReleases()
    }

  };
</script>

