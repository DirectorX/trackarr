<template>
  <v-app>
    <v-app-bar v-if="this.$vuetify.breakpoint.name == 'xs'" app color="accent-4" dark>
      <v-menu>
        <template v-slot:activator="{ on }">
          <v-btn icon v-on="on">
            <v-icon>mdi-menu</v-icon>
          </v-btn>
        </template>

        <v-list>
          <v-list-item link to="/">

            <v-list-item-title>
              Home
            </v-list-item-title>

          </v-list-item>
          <v-list-item link to="/logs">

            <v-list-item-title>
              Logs
            </v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
      <v-app-bar-title>Trackarr</v-app-bar-title>
    </v-app-bar>
    <v-navigation-drawer v-if="this.$vuetify.breakpoint.name != 'xs'" app mini-variant permanent>
      <v-list dense nav>
        <v-list-item-title class="pa-3 text-center text-uppercase">
          <span class="font-weight-black">T</span>
        </v-list-item-title>
        <v-divider></v-divider>
        <v-list-item link to="/">
          <v-list-item-icon>
            <v-icon>mdi-home</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
          </v-list-item-content>

        </v-list-item>
        <v-list-item link to="/logs">
          <v-list-item-icon>
            <v-icon>mdi-folder</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>
    <v-content>
      <transition name="fade" mode="out-in">
        <router-view />
      </transition>

    </v-content>
  </v-app>
</template>

<script>
  require('@/assets/styles/trackarr.css')
  export default {
    name: 'App',
    data: () => ({

    })
  };
</script>