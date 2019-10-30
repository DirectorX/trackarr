<template>
  <v-app>
    <v-navigation-drawer app mini-variant permanent>
      <v-list dense nav>
        <v-list-item-title class="pa-3 text-center text-uppercase">
          <span class="font-weight-black">T</span>
        </v-list-item-title>
        <v-divider></v-divider>
        <v-list-item link>
          <v-list-item-icon>
            <v-icon>mdi-home</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
          </v-list-item-content>
       
        </v-list-item>
        <v-list-item link>
          <v-list-item-icon>
            <v-icon>mdi-folder</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
          </v-list-item-content>
        </v-list-item>
        <v-list-item link>
          <v-list-item-icon>
            <v-icon>mdi-cogs</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>
    <v-content>
      <router-view />
    </v-content>
  </v-app>
</template>

<script>
  export default {
    name: 'App',
    data: () => ({
      //
    }),
  };
</script>