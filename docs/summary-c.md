# Detail Page info

```c
/*
 * Manage the state icon in the panel and the associated tooltip
 * Composes the expanded tooltip message
 */
static gboolean gapc_monitor_update_tooltip_msg(PGAPC_MONITOR pm)
{
   gchar *pchx = NULL, *pmsg = NULL, *ptitle = NULL, *pch5 = NULL; 
   gchar *pch5a = NULL, *pch5b = NULL;
   gchar *pch1 = NULL, *pch2 = NULL, *pch3 = NULL, *pch4 = NULL;
   gchar *pch6 = NULL, *pch7 = NULL, *pch8 = NULL, *pch9 = NULL;
   gchar *pchb = NULL, *pchc = NULL, *pchd = NULL, *pche = NULL;
   gchar *pcha = NULL, *pmview = NULL, *pchw = NULL;
   GtkWidget *w = NULL;
   gdouble d_value = 0.0, d_loadpct = 0.0, d_watt = 0.0;
   gboolean b_flag = FALSE, b_valid = FALSE;
   gint i_series = 0;
   GtkTreeIter miter;
   GdkPixbuf *pixbuf;

   g_return_val_if_fail(pm != NULL, TRUE);

   if (pm->b_run != TRUE)
      return TRUE;

   w = g_hash_table_lookup(pm->pht_Widgets, "StatusBar");

   pm->i_icon_index = GAPC_ICON_ONLINE;

   pch1 = g_hash_table_lookup(pm->pht_Status, "UPSNAME");
   pch2 = g_hash_table_lookup(pm->pht_Status, "HOSTNAME");
   if (pch2 == NULL) {
      pch2 = pm->pch_host;
   }
   if (pch2 == NULL) {
      pch2 = "unknown";
   }
   pch3 = g_hash_table_lookup(pm->pht_Status, "STATUS");
   if (pch3 == NULL) {
      pch3 = "NISERR";
   }
   pch4 = g_hash_table_lookup(pm->pht_Status, "NUMXFERS");
   pch5a = g_hash_table_lookup(pm->pht_Status, "TONBATT");
   if (pch5a == NULL) {
      pch5a = " ";
   }
   pch5b = g_hash_table_lookup(pm->pht_Status, "CUMONBATT");
   if (pch5b == NULL) {
      pch5b = " ";
   }
   pch5 = g_hash_table_lookup(pm->pht_Status, "XONBATT");
   if (pch5 == NULL) {
      pch5 = " ";
   }
   pch6 = g_hash_table_lookup(pm->pht_Status, "LINEV");
   pch7 = g_hash_table_lookup(pm->pht_Status, "BCHARGE");
   if (pch7 == NULL) {
      pch7 = "n/a";
   }
   pch8 = g_hash_table_lookup(pm->pht_Status, "LOADPCT");
   if (pch8 != NULL) {
       d_loadpct = g_strtod (pch8, NULL);
       d_loadpct /= 100.0;
   }
   pchw = g_hash_table_lookup(pm->pht_Status, "NOMPOWER");
   if (pchw == NULL) {
       d_watt = d_loadpct * pm->i_watt;        
   } else {
       pm->i_watt = g_strtod (pchw, NULL);
       d_watt = d_loadpct * pm->i_watt;        
   }
   pch9 = g_hash_table_lookup(pm->pht_Status, "TIMELEFT");
   pcha = g_hash_table_lookup(pm->pht_Status, "VERSION");
   pchb = g_hash_table_lookup(pm->pht_Status, "STARTTIME");
   pchc = g_hash_table_lookup(pm->pht_Status, "MODEL");
   pchd = g_hash_table_lookup(pm->pht_Status, "UPSMODE");
   pche = g_hash_table_lookup(pm->pht_Status, "CABLE");

   if (pm->b_data_available) {
      d_value = g_strtod(pch7, NULL);
      pchx = NULL;
      if (g_strrstr(pch3, "COMMLOST") != NULL) {
         pchx = " cable un-plugged...";
         pm->i_icon_index = GAPC_ICON_UNPLUGGED;
         b_flag = TRUE;
      } else if ((d_value < 99.0) && (g_strrstr(pch3, "LINE") != NULL)) {
         pchx = " and charging...";
         pch3 = "CHARGING";
         pm->i_icon_index = GAPC_ICON_CHARGING;
      } else if (g_strrstr(pch3, "BATT") != NULL) {
         pchx = " on battery...";
         pm->i_icon_index = GAPC_ICON_ONBATT;
      }
   } else {
      b_flag = TRUE;
      pchx = " NIS network error...";
      pch3 = "NISERR";
      g_hash_table_replace(pm->pht_Status, g_strdup("STATUS"), g_strdup(pch3));
      pm->i_icon_index = GAPC_ICON_NETWORKERROR;
      for (i_series = 0; i_series < pm->phs.plg->i_num_series; i_series++) {
         gapc_util_point_filter_set(&(pm->phs.sq[i_series]), 0.0);
      }
   }

   if (b_flag) {
      ptitle = g_strdup_printf("<span foreground=\"red\" size=\"large\">"
         "%s@%s\nis %s%s" "</span>",
         (pch1 != NULL) ? pch1 : "unknown",
         (pch2 != NULL) ? pch2 : "unknown",
         (pch3 != NULL) ? pch3 : "n/a", (pchx != NULL) ? pchx : " ");
   } else {
      ptitle = g_strdup_printf("<span foreground=\"blue\" size=\"large\">"
         "%s@%s\nis %s%s" "</span>",
         (pch1 != NULL) ? pch1 : "unknown",
         (pch2 != NULL) ? pch2 : "unknown",
         (pch3 != NULL) ? pch3 : "n/a", (pchx != NULL) ? pchx : " ");
   }

   pmsg = g_strdup_printf("%s@%s\nStatus: %s%s\n"
      "Refresh occurs every %3.1f seconds\n"
      "----------------------------------------------------------\n"
      "%s Outage[s]\n" 
      "Last on %s\n" 
      "%s Utility VAC\n"
      "%s Battery Charge\n" 
      "%s UPS Load\n"
      "%3.1f of %d watts\n" 
      "%s Remaining\n"
      "----------------------------------------------------------\n"
      "Build: %s\n" 
      "Started: %s\n"
      "----------------------------------------------------------\n"
      "Model: %s\n" 
      " Mode: %s\n" 
      "Cable: %s",
      (pch1 != NULL) ? pch1 : "unknown",
      (pch2 != NULL) ? pch2 : "unknown",
      (pch3 != NULL) ? pch3 : "n/a",
      (pchx != NULL) ? pchx : " ",
      pm->d_refresh,
      (pch4 != NULL) ? pch4 : "n/a",
      (pch5 != NULL) ? pch5 : "n/a",
      (pch6 != NULL) ? pch6 : "n/a",
      (pch7 != NULL) ? pch7 : "n/a",
      (pch8 != NULL) ? pch8 : "n/a", 
      d_watt, pm->i_watt,
      (pch9 != NULL) ? pch9 : "n/a",
      (pcha != NULL) ? pcha : "n/a",
      (pchb != NULL) ? pchb : "n/a",
      (pchc != NULL) ? pchc : "n/a",
      (pchd != NULL) ? pchd : "n/a", (pche != NULL) ? pche : "n/a");


   switch (pm->i_icon_index) {
   case GAPC_ICON_NETWORKERROR:
      pmview = g_strdup_printf("<span foreground=\"red\" size=\"large\">"
         "<b><i>%s@%s</i></b></span>\n"
         "NIS network connection not Responding!",
         (pch1 != NULL) ? pch1 : "unknown", (pch2 != NULL) ? pch2 : "unknown");
      break;
   case GAPC_ICON_UNPLUGGED:
      pmview = g_strdup_printf("<span foreground=\"red\" size=\"large\">"
         "<b><i>%s@%s</i></b></span>\n"
         "%s",
         (pch1 != NULL) ? pch1 : "unknown",
         (pch2 != NULL) ? pch2 : "unknown", (pchx != NULL) ? pchx : " un-plugged");
      break;
   case GAPC_ICON_CHARGING:
      pmview = g_strdup_printf("<span foreground=\"blue\">"
         "<b><i>%s@%s</i></b></span>\n"
         "%s Outage, Last on %s\n"
         "%s VAC, %s Charge, %3.1f of %d watts\n"
         "%s Remaining, %s total on battery",
         (pch1 != NULL) ? pch1 : "unknown",
         (pch2 != NULL) ? pch2 : "unknown",
         (pch4 != NULL) ? pch4 : "n/a",
         (pch5 != NULL) ? pch5 : "n/a",
         (pch6 != NULL) ? pch6 : "n/a",
         (pch7 != NULL) ? pch7 : "n/a", d_watt, pm->i_watt,
         (pch9 != NULL) ? pch9 : "n/a", (pch5b != NULL) ? pch5b : "n/a");
      break;
   case GAPC_ICON_ONBATT:
      pmview = g_strdup_printf("<span foreground=\"yellow\">"
         "<b><i>%s@%s</i></b></span>\n"
         "%s Outage, Last on %s\n"
         "%s Charge, %s total on battery\n"
         "%s Remaining, %s on battery\n"
         "%3.1f of %d watts",
         (pch1 != NULL) ? pch1 : "unknown",
         (pch2 != NULL) ? pch2 : "unknown",
         (pch4 != NULL) ? pch4 : "n/a",
         (pch5 != NULL) ? pch5 : "n/a",
         (pch7 != NULL) ? pch7 : "n/a", 
         (pch5b != NULL) ? pch5b : "n/a",
         (pch9 != NULL) ? pch9 : "n/a", 
         (pch5a != NULL) ? pch5a : "n/a ",
         d_watt, pm->i_watt);
      break;
   case GAPC_ICON_ONLINE:
   case GAPC_ICON_DEFAULT:
   default:
      pmview = g_strdup_printf("<b><i>%s@%s</i></b>\n"
         "%s Outage, Last on %s\n"
         "%s VAC, %s Charge %s, %3.0f of %d watts",
         (pch1 != NULL) ? pch1 : "unknown",
         (pch2 != NULL) ? pch2 : "unknown",
         (pch4 != NULL) ? pch4 : "n/a",
         (pch5 != NULL) ? pch5 : "n/a",
         (pch6 != NULL) ? pch6 : "n/a",
         (pch7 != NULL) ? pch7 : "n/a", 
         (pchx != NULL) ? pchx : " ",
         d_watt, pm->i_watt);
      break;
   }

   pixbuf = pm->my_icons[pm->i_icon_index];
   if (pm->i_old_icon_index != pm->i_icon_index) {
      b_flag = TRUE;
   } else {
      b_flag = FALSE;
   }

   if ((pm->tooltips != NULL) && (pm->tray_icon != NULL)) {
      gtk_tooltips_set_tip(pm->tooltips, GTK_WIDGET(pm->tray_icon), pmsg, NULL);
      if (b_flag) {
         gapc_util_change_icons(pm);
      }
   }

   b_valid =
      gapc_util_treeview_get_iter_from_monitor(pm->monitor_model, &miter,
      pm->cb_monitor_num);
   if (b_valid) {
      if (b_flag) {
         gtk_list_store_set(GTK_LIST_STORE(pm->monitor_model), &miter,
            GAPC_MON_STATUS, pmview, GAPC_MON_UPSSTATE, pch3,
            GAPC_MON_ICON, pixbuf, -1);
      } else {
         gtk_list_store_set(GTK_LIST_STORE(pm->monitor_model), &miter,
            GAPC_MON_STATUS, pmview, GAPC_MON_UPSSTATE, pch3, -1);
      }
   }

   if ((w = g_hash_table_lookup(pm->pht_Widgets, "TitleStatus"))) {
      gtk_label_set_markup(GTK_LABEL(w), ptitle);
      lg_graph_set_chart_title (pm->phs.plg, ptitle);
      g_snprintf(pm->ch_title_info, GAPC_MAX_TEXT, "%s", ptitle);
      
      lg_graph_draw ( pm->phs.plg );
   }

   g_free(pmsg);
   g_free(ptitle);

/*  g_free (pmview); */

   return b_flag;
}

/*
 * main data updating routine.
 * -- collects and pushes data to all ui
 */
static gint gapc_monitor_update(PGAPC_MONITOR pm)
{
   gint i_x = 0;
   GtkWidget *win = NULL, *w = NULL;
   gchar *pch  = NULL, *pch1 = NULL, *pch2 = NULL, *pch3 = NULL;
   gchar *pch4 = NULL, *pch5 = NULL, *pch6 = NULL;
   gdouble dValue = 0.00, dScale = 0.0, dtmp = 0.0, dCharge = 0.0, d_loadpct = 0.0;
   gchar ch_buffer[GAPC_MAX_TEXT];
   PGAPC_BAR_H pbar = NULL;

   g_return_val_if_fail(pm != NULL, FALSE);

   if (pm->window == NULL)         /* not created yet */
      return TRUE;

   if (pm->b_run == FALSE)
      return FALSE;

   if (pm->b_data_available == FALSE)
      return FALSE;

   w = g_hash_table_lookup(pm->pht_Widgets, "StatusPage");
   if (gapc_util_text_view_clear_buffer(GTK_WIDGET(w))) {
      return FALSE;
   }
   for (i_x = 1; pm->pach_status[i_x] != NULL; i_x++) {
      gapc_util_text_view_append(GTK_WIDGET(w), pm->pach_status[i_x]);
   }

   w = g_hash_table_lookup(pm->pht_Widgets, "EventsPage");
   gapc_util_text_view_clear_buffer(GTK_WIDGET(w));
   for (i_x = 0; pm->pach_events[i_x] != NULL; i_x++) {
      gapc_util_text_view_prepend(GTK_WIDGET(w), pm->pach_events[i_x]);
   }

   /*
    *  compute graphic points */
   pch = g_hash_table_lookup(pm->pht_Status, "LINEV");
   if (pch == NULL) {
      pch = "n/a";
   }
   dValue = g_strtod(pch, NULL);
   dScale = (( dValue - 200 ) > 1) ? 230.0 : 120.0;
   dValue /= dScale;
   gapc_util_point_filter_set(&(pm->phs.sq[0]), dValue);
   pbar = g_hash_table_lookup(pm->pht_Status, "HBar1");
   pbar->d_value = dValue;
   g_snprintf(pbar->c_text, sizeof(pbar->c_text), "%s from Utility", pch);
   w = g_hash_table_lookup(pm->pht_Widgets, "HBar1-Widget");
   if (GTK_WIDGET_DRAWABLE(w))
      gdk_window_invalidate_rect(w->window, &pbar->rect, FALSE);

   pch = g_hash_table_lookup(pm->pht_Status, "BATTV");
   if (pch == NULL) {
      pch = "n/a";
   }
   pch1 = g_hash_table_lookup(pm->pht_Status, "NOMBATTV");
   if (pch1 == NULL) {
      pch1 = "n/a";
   }
   dValue = g_strtod(pch, NULL);
   dScale = g_strtod(pch1, NULL);
   if (dScale == 0.0)
      dScale = ((gint) (dValue - 20)) ? 24 : 12;
   dValue /= dScale;
   gapc_util_point_filter_set(&(pm->phs.sq[4]), dValue);
   pbar = g_hash_table_lookup(pm->pht_Status, "HBar2");
   pbar->d_value = (dValue > 1.0) ? 1.0 : dValue;
   g_snprintf(pbar->c_text, sizeof(pbar->c_text), "%s DC on Battery", pch);

   w = g_hash_table_lookup(pm->pht_Widgets, "HBar2-Widget");
   if (GTK_WIDGET_DRAWABLE(w))
      gdk_window_invalidate_rect(w->window, &pbar->rect, FALSE);

   pch = g_hash_table_lookup(pm->pht_Status, "BCHARGE");
   if (pch == NULL) {
      pch = "n/a";
   }
   dCharge = dValue = g_strtod(pch, NULL);
   dValue /= 100.0;
   gapc_util_point_filter_set(&(pm->phs.sq[3]), dValue);
   pbar = g_hash_table_lookup(pm->pht_Status, "HBar3");
   pbar->d_value = dValue;
   g_snprintf(pbar->c_text, sizeof(pbar->c_text), "%s Battery Charge", pch);
   w = g_hash_table_lookup(pm->pht_Widgets, "HBar3-Widget");
   if (GTK_WIDGET_DRAWABLE(w))
      gdk_window_invalidate_rect(w->window, &pbar->rect, FALSE);

   pch = g_hash_table_lookup(pm->pht_Status, "LOADPCT");
   if (pch == NULL) {
      pch = "n/a";
   }
   dValue = g_strtod(pch, NULL);
   dtmp = dValue /= 100.0;
   d_loadpct = dtmp;
   gapc_util_point_filter_set(&(pm->phs.sq[1]), dValue);
   pbar = g_hash_table_lookup(pm->pht_Status, "HBar4");
   pbar->d_value = (dValue > 1.0) ? 1.0 : dValue;
   g_snprintf(pbar->c_text, sizeof(pbar->c_text), "%s", pch);

   w = g_hash_table_lookup(pm->pht_Widgets, "HBar4-Widget");
   if (GTK_WIDGET_DRAWABLE(w))
      gdk_window_invalidate_rect(w->window, &pbar->rect, FALSE);

   pch = g_hash_table_lookup(pm->pht_Status, "TIMELEFT");
   if (pch == NULL) {
      pch = "n/a";
   }
   dValue = g_strtod(pch, NULL);
   dScale = dValue / (1 - dtmp);
   dValue /= dScale;
   gapc_util_point_filter_set(&(pm->phs.sq[2]), dValue);
   pbar = g_hash_table_lookup(pm->pht_Status, "HBar5");
   pbar->d_value = dValue;
   g_snprintf(pbar->c_text, sizeof(pbar->c_text), "%s Remaining", pch);
   w = g_hash_table_lookup(pm->pht_Widgets, "HBar5-Widget");
   if (GTK_WIDGET_DRAWABLE(w))
      gdk_window_invalidate_rect(w->window, &pbar->rect, FALSE);

   /*
    * information window update */
   win = g_hash_table_lookup(pm->pht_Widgets, "SoftwareInformation");
   pch = g_hash_table_lookup(pm->pht_Status, "VERSION");
   pch1 = g_hash_table_lookup(pm->pht_Status, "UPSNAME");
   pch2 = g_hash_table_lookup(pm->pht_Status, "CABLE");
   pch3 = g_hash_table_lookup(pm->pht_Status, "UPSMODE");
   pch4 = g_hash_table_lookup(pm->pht_Status, "STARTTIME");
   pch5 = g_hash_table_lookup(pm->pht_Status, "STATUS");   
   g_snprintf(ch_buffer, sizeof(ch_buffer),
      "<span foreground=\"blue\">" "%s\n%s\n%s\n%s\n%s\n%s" "</span>",
      (pch != NULL) ? pch : "N/A", (pch1 != NULL) ? pch1 : "N/A",
      (pch2 != NULL) ? pch2 : "N/A", (pch3 != NULL) ? pch3 : "N/A",
      (pch4 != NULL) ? pch4 : "N/A", (pch5 != NULL) ? pch5 : "N/A");
   gtk_label_set_markup(GTK_LABEL(win), ch_buffer);

   win = g_hash_table_lookup(pm->pht_Widgets, "PerformanceSummary");
   pch = g_hash_table_lookup(pm->pht_Status, "SELFTEST");
   pch1 = g_hash_table_lookup(pm->pht_Status, "NUMXFERS");
   pch2 = g_hash_table_lookup(pm->pht_Status, "LASTXFER");
   pch3 = g_hash_table_lookup(pm->pht_Status, "XONBATT");
   pch4 = g_hash_table_lookup(pm->pht_Status, "XOFFBATT");
   pch5 = g_hash_table_lookup(pm->pht_Status, "TONBATT");
   pch6 = g_hash_table_lookup(pm->pht_Status, "CUMONBATT");
   g_snprintf(ch_buffer, sizeof(ch_buffer),
      "<span foreground=\"blue\">" "%s\n%s\n%s\n%s\n%s\n%s\n%s" "</span>",
      (pch != NULL) ? pch : "N/A", (pch1 != NULL) ? pch1 : "N/A",
      (pch2 != NULL) ? pch2 : "N/A", (pch3 != NULL) ? pch3 : "N/A",
      (pch4 != NULL) ? pch4 : "N/A", (pch5 != NULL) ? pch5 : "N/A",
      (pch6 != NULL) ? pch6 : "N/A");
   gtk_label_set_markup(GTK_LABEL(win), ch_buffer);

   win = g_hash_table_lookup(pm->pht_Widgets, "ProductInformation");
   pch = g_hash_table_lookup(pm->pht_Status, "MODEL");
   pch1 = g_hash_table_lookup(pm->pht_Status, "SERIALNO");
   pch2 = g_hash_table_lookup(pm->pht_Status, "MANDATE");
   pch3 = g_hash_table_lookup(pm->pht_Status, "FIRMWARE");
   pch4 = g_hash_table_lookup(pm->pht_Status, "BATTDATE");
   pch6 = g_hash_table_lookup(pm->pht_Status, "NOMPOWER");
   if (pch6 == NULL) {
       dValue =  d_loadpct * pm->i_watt; 
   } else {
       pm->i_watt = g_strtod (pch6, NULL);
       dValue =  d_loadpct * pm->i_watt; 
   }
   g_snprintf(ch_buffer, sizeof(ch_buffer),
      "<span foreground=\"blue\">" "%s\n%s\n%s\n%s\n%s\n%3.1f of %d" "</span>",
      (pch != NULL) ? pch : "N/A", (pch1 != NULL) ? pch1 : "N/A",
      (pch2 != NULL) ? pch2 : "N/A", (pch3 != NULL) ? pch3 : "N/A",
      (pch4 != NULL) ? pch4 : "N/A", dValue, pm->i_watt);
   gtk_label_set_markup(GTK_LABEL(win), ch_buffer);

   return TRUE;
}

```

