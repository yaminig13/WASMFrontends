#include <stdlib.h>
#include <stdarg.h>
#include <emscripten.h>
#include <mruby.h>
#include <mruby/irep.h>
#include <mruby/array.h>
#include <mruby/proc.h>
#include <mruby/compile.h>
#include <mruby/dump.h>
#include <mruby/string.h>
#include <mruby/variable.h>
#include <mruby/throw.h>

mrb_value app;
mrb_state *mrb;
mrbc_context *c;


mrb_value
add_event_listener(mrb_state *mrb, mrb_value self){
  mrb_value selector, event, id;
  mrb_get_args(mrb, "SSS", &selector, &event, &id);

  EM_ASM_({
    var selector = UTF8ToString($0);
    var eventName = UTF8ToString($1);
    var id = UTF8ToString($2);
    var elements;

    if (selector === 'document') {
      elements = [window.document];
    } else if (selector === 'body') {
      elements = [window.document.body];
    } else {
      elements = document.querySelectorAll(selector);
    }

    for (var i = 0; i < elements.length; i++) {
      var element = elements[i];

      element.addEventListener(
          eventName, function(event)
          {
            Module.ccall(
                'event',
                'void',
                [ 'string', 'string', 'string' ],
                [ stringifyEvent(event), id ]);

            render();
          });
    };
  }, RSTRING_PTR(selector), RSTRING_PTR(event), RSTRING_PTR(id));
  return mrb_nil_value();
}

mrb_value
get_location(mrb_state *mrb, mrb_value self)
{
  mrb_value id;
  mrb_get_args(mrb, "S", &id);

  EM_ASM_(
      {
        var id = UTF8ToString($0);

        navigator.geolocation.getCurrentPosition(function(event)
                                                 {
                                                   event = {
                                                     latitude : event.coords.latitude,
                                                     longitude : event.coords.longitude
                                                   };
                                                   Module.ccall(
                                                       'display_location',
                                                       'void',
                                                       [ 'string', 'string' ],
                                                       [ stringifyEvent(event), id ]);
                                                   render();
                                                 });
      },
      RSTRING_PTR(id));
  return mrb_nil_value();
}

mrb_value get_value(mrb_state *mrb, mrb_value self)
{
  mrb_value selector, id;
  mrb_get_args(mrb, "SS", &selector, &id);

  EM_ASM_(
      {
        const id = UTF8ToString($1);
        const selector = UTF8ToString($0);

        const data = document.querySelector(selector).value;
        Module.ccall('display_value',
                     'void',
                     [ 'string', 'string' ],
                     [ data, id ]);
        render();
      },
      RSTRING_PTR(selector), RSTRING_PTR(id));
  return mrb_nil_value();
}

// mrb_value get_element(mrb_state *mrb, mrb_value self)
// {
//   mrb_value selector, id;
//   mrb_get_args(mrb, "SS", &selector, &id);

//   EM_ASM_(
//       {
//         const id = UTF8ToString($1);
//         const selector = UTF8ToString($0);

//         var element = document.querySelector(selector);
//         Module.ccall('display_value',
//                      'void',
//                      [ 'string', 'string' ],
//                      [ element, id ]);
//         render();
//       },
//       RSTRING_PTR(selector), RSTRING_PTR(id));
//   return mrb_nil_value();
// }

mrb_value
http_request(mrb_state *mrb, mrb_value self){
  mrb_value url, id;
  mrb_get_args(mrb, "SS", &url, &id);

  EM_ASM_({
    const response = fetch(UTF8ToString($0));
    const id = UTF8ToString($1);

    response.then(r => r.text()).then(text => {
      Module.ccall('http_response',
        'void',
        ['string', 'string'],
        [JSON.stringify({body: text}), id]
      );

      render();
    });
  }, RSTRING_PTR(url), RSTRING_PTR(id));
  return mrb_nil_value();
}

mrb_value http_post(mrb_state *mrb, mrb_value self)
{
  mrb_value url, data;
  mrb_get_args(mrb, "SS", &url, &data);

  EM_ASM_(
      {
        fetch(UTF8ToString($0), {
          method : "POST",
          headers : {
            'content-type' : 'application/json; charset=utf-8',
          },
          body : JSON.stringify(
              {
                image64 : UTF8ToString($1)
              }),
        })
            .then(() => {alert("Review Submitted!")});
      },
      RSTRING_PTR(url), RSTRING_PTR(data));
  return mrb_nil_value();
}

mrb_value get_session(mrb_state *mrb, mrb_value self)
{
  mrb_value sessionName, id;
  mrb_get_args(mrb, "SS", &sessionName, &id);

  EM_ASM_(
      {
        const id = UTF8ToString($1);
        const sessionResp = sessionStorage.getItem(UTF8ToString($0));
        Module.ccall('session_value',
                     'void',
                     [ 'string', 'string' ],
                     [ sessionResp, id ]);
        render();
      },
      RSTRING_PTR(sessionName), RSTRING_PTR(id));
  return mrb_nil_value();
}

mrb_value get_cookie(mrb_state *mrb, mrb_value self)
{
  mrb_value id;
  mrb_get_args(mrb, "S", &id);
  EM_ASM_(
      {
        const id = UTF8ToString($0);

        const cookieValues = document.cookie;

        Module.ccall('cookie_value',
                     'void',
                     [ 'string', 'string' ],
                     [ cookieValues, id ]);
        render();
      },
      RSTRING_PTR(id));

  return mrb_nil_value();
}

mrb_value set_cookie(mrb_state *mrb, mrb_value self)
{
  mrb_value cookieName;
  mrb_get_args(mrb, "S", &cookieName);

  EM_ASM_(
      {
        document.cookie = UTF8ToString($0);
      },
      RSTRING_PTR(cookieName));
  return mrb_nil_value();
}

mrb_value set_session(mrb_state *mrb, mrb_value self)
{
  mrb_value sessionName, sessionValue;
  mrb_get_args(mrb, "SS", &sessionName, &sessionValue);

  EM_ASM_(
      {
        sessionStorage.setItem(UTF8ToString($0), UTF8ToString($1));
      },
      RSTRING_PTR(sessionName), RSTRING_PTR(sessionValue));
  return mrb_nil_value();
}

mrb_value set_local(mrb_state *mrb, mrb_value self)
{
  mrb_value localName, localValue;
  mrb_get_args(mrb, "SS", &localName, &localValue);

  EM_ASM_(
      {
        localStorage.setItem(UTF8ToString($0), UTF8ToString($1));
      },
      RSTRING_PTR(localName), RSTRING_PTR(localValue));
  return mrb_nil_value();
}
mrb_value get_local(mrb_state *mrb, mrb_value self)
{
  mrb_value localName, id;
  mrb_get_args(mrb, "SS", &localName, &id);

  EM_ASM_(
      {
        const id = UTF8ToString($1);
        const localResp = localStorage.getItem(UTF8ToString($0));
        Module.ccall('local_value', 'void', [ 'string', 'string' ], [ localResp, id ]);
        render();
      },
      RSTRING_PTR(localName), RSTRING_PTR(id));
  return mrb_nil_value();
}

mrb_value open_connection(mrb_state *mrb, mrb_value self)
{
  mrb_value url, clientMsg, id;
  mrb_get_args(mrb, "SSS", &url, &clientMsg, &id);

  EM_ASM_(
      {
        const id = UTF8ToString($2);
        connection = new WebSocket(UTF8ToString($0));

        connection.addEventListener(
            'open', (event) => {
              connection.send(UTF8ToString($1));
            });

        connection.onmessage = (e) =>
        {
          Module.ccall('send_msg',
                       'void',
                       [ 'string', 'string' ],
                       [ e.data, id ]);
          render();
        }
      },
      RSTRING_PTR(url), RSTRING_PTR(clientMsg), RSTRING_PTR(id));
  return mrb_nil_value();
}
mrb_value navigate_to(mrb_state *mrb, mrb_value self)
{
  mrb_value url, customId;
  mrb_get_args(mrb, "SS", &url, &customId);

  EM_ASM_(
      {
        const url = UTF8ToString($0);
        const customId = UTF8ToString($1);

        window.location = (url + customId)
      },
      RSTRING_PTR(url), RSTRING_PTR(customId));
  return mrb_nil_value();
}

mrb_value get_param(mrb_state *mrb, mrb_value self)
{
  mrb_value pname, id;
  mrb_get_args(mrb, "SS", &pname, &id);

  EM_ASM_(
      {
        const id = UTF8ToString($1);

        const url_string = window.location.href;
        const url = new URL(url_string);
        const pvalue = url.searchParams.get(UTF8ToString($0));
        Module.ccall('param_value', 'void', [ 'string', 'string' ], [ pvalue, id ]);
        // render();
      },
      RSTRING_PTR(pname), RSTRING_PTR(id));
  return mrb_nil_value();
}
mrb_value vibrate_for(mrb_state *mrb, mrb_value self)
{
  mrb_value time;
  mrb_get_args(mrb, "i", &time);

  EM_ASM_(
      {
        const time = UTF8ToString($0);

        vibrated = window.navigator.vibrate(time);
        console.log("vibrated:" + vibrated);
      });
  return mrb_nil_value();
}

mrb_value send_notif(mrb_state *mrb, mrb_value self)
{
  mrb_value t, b;
  mrb_get_args(mrb, "SS", &t, &b);

  EM_ASM_(
      {
        if (window.Notification &&Notification.permission !== "denied")
        {
          const title = UTF8ToString($0);
          Notification.requestPermission(function(status) { var notif = new Notification(title, {
                                                              body : UTF8ToString($1)
                                                            }); });
        }
      },
      RSTRING_PTR(t), RSTRING_PTR(b));
  return mrb_nil_value();
}

mrb_value show_alert(mrb_state *mrb, mrb_value self)
{
  mrb_value msg;
  mrb_get_args(mrb, "S", &msg);

  EM_ASM_(
      {alert(UTF8ToString($0))},
      RSTRING_PTR(msg));
  return mrb_nil_value();
}

mrb_value receive_msg(mrb_state *mrb, mrb_value self)
{
  mrb_value worker, id;
  mrb_get_args(mrb, "SS", &worker, &id);

  EM_ASM_(
      {
        id = UTF8ToString($1);
        activityWorker = JSON.parse(UTF8ToString($0));
        console.log("outside");
        activityWorker.onmessage = function(e)
        {
          console.log("inside");
          Module.ccall('worker_fn', 'void', [ 'string', 'string' ], [ stringifyEvent(e), id ]);
        }
      },
      RSTRING_PTR(worker), RSTRING_PTR(id));
  return mrb_nil_value();
}

mrb_value start_worker(mrb_state *mrb, mrb_value self)
{
  mrb_value url, id;
  mrb_get_args(mrb, "SS", &url, &id);

  EM_ASM_(
      {
        id = UTF8ToString($1);
        activityWorker = new Worker(UTF8ToString($0));
        // activityWorker.onmessage = function(e)
        // {
        //   Module.ccall('worker_fn', 'void', [ 'string', 'string' ], [ stringifyEvent(e), id ]);
        // }
        // console.log(typeof stringifyEvent(activityWorker));

        Module.ccall('worker_fn', 'void', [ 'string', 'string' ], [ stringifyEvent(activityWorker), id ]);
      },
      RSTRING_PTR(url), RSTRING_PTR(id));
  return mrb_nil_value();
}

mrb_value enable_camera(mrb_state *mrb, mrb_value self)
{
  mrb_value id;
  mrb_get_args(mrb, "S", &id);
  EM_ASM_(
      {
        id = UTF8ToString($0);

        var video = document.querySelector("video");
        var canvas = document.querySelector("canvas");
        var context = canvas.getContext("2d");

        document.querySelector(".snap").addEventListener(
            "click", function()
            {
              context.drawImage(video, 0, 0, 200, 200);
              var imageData = canvas.toDataURL();

              document.querySelector(".cameraInput")
                  .setAttribute("value", imageData);
              Module.ccall('camera_output', 'void', [ 'string', 'string' ], [ imageData, id ]);
              render();

            });
        if (navigator.mediaDevices && navigator.mediaDevices.getUserMedia)
        {
          navigator.mediaDevices.getUserMedia({video : true}).then(function(stream)
                                                                   {
                                                                     video.srcObject = stream;
                                                                     video.play();
                                                                   })
        }
      },
      RSTRING_PTR(id));
  return mrb_nil_value();
}

mrb_value send_pushNotif(mrb_state *mrb, mrb_value self)
{
  EM_ASM_(
      {
        navigator.serviceWorker.register("serviceWorker.js");

        navigator.serviceWorker.ready.then(function(registration){
                                               return registration.pushManager.getSubscription().then(async function(subscription)
                                                                                                      {
                                                                                                        if (subscription)
                                                                                                        {
                                                                                                          return subscription;
                                                                                                        }

                                                                                                        const response = await fetch("http://localhost:8080/vapidPublicKey");
                                                                                                        var vapidPublicKey = await response.text();
                                                                                                        const convertedVapidKey = urlBase64ToUint8Array(vapidPublicKey);
                                                                                                        return registration.pushManager.subscribe({
                                                                                                          userVisibleOnly : true,
                                                                                                          applicationServerKey : convertedVapidKey
                                                                                                        });
                                                                                                      })})
            .then(function(subscription)
                  {
                    fetch("http://localhost:8080/register", {
                      method : "post",
                      headers : {
                        "Content-type" : "application/json"
                      },
                      body : JSON.stringify({
                        subscription : subscription
                      }),
                    });

                    fetch("http://localhost:8080/sendNotification", {
                      method : "post",
                      headers : {
                        "Content-type" : "application/json"
                      },
                      body : JSON.stringify({
                        subscription : subscription,
                        delay : 120
                      }),
                    });
                  });
      });
  return mrb_nil_value();
}

mrb_value create_db(mrb_state *mrb, mrb_value self)
{
  mrb_value name, items;
  mrb_get_args(mrb, "SS", &name, &items);

  EM_ASM_(
      {
        const dbName = UTF8ToString($0);
        var products = JSON.parse(UTF8ToString($1));

        // console.log(products);
        var request = indexedDB.open(dbName, 1);
        request.onerror = function(event)
        {
          console.error("error found");
        };

        request.onupgradeneeded = function(event)
        {
          var db = event.target.result;

          var objectStore = db.createObjectStore("products", {keyPath : "id"});
          objectStore.createIndex("title", "title", {unique : false});
          objectStore.transaction.oncomplete = function(event)
          {
            var productStore = db.transaction("products", "readwrite").objectStore("products");
            products.forEach(function(product) {
              productStore.add(product);
            });
          };
        };
      },
      RSTRING_PTR(name), RSTRING_PTR(items));
  return mrb_nil_value();
}

mrb_value open_db(mrb_state *mrb, mrb_value self)
{
  mrb_value name, term, id;
  mrb_get_args(mrb, "SSS", &name, &term, &id);

  EM_ASM_(
      {
        const dbName = UTF8ToString($0);
        var searchTerm = UTF8ToString($1);
        const id = UTF8ToString($2);

        searchTerm = searchTerm.toLowerCase();
        var request = indexedDB.open(dbName, 1);
        request.onerror = function(event)
        {
          console.error("error found");
        };

        request.onsuccess = function(event)
        {
          console.log("connection open");
          var db = event.target.result;
          var txn = db.transaction("products", "readwrite");
          let tempProducts = txn.objectStore("products");
          let titleIndex = tempProducts.index("title");
          let query = titleIndex.getAll(searchTerm);

          query.onsuccess = function()
          {
            Module.ccall('indexeddb_response', 'void', [ 'string', 'string' ], [ JSON.stringify(query.result), id ]);
          };
        };
      },
      RSTRING_PTR(name), RSTRING_PTR(term), RSTRING_PTR(id));
  return mrb_nil_value();
}

int main(int argc, const char *argv[])
{
  struct RClass *dom_class, *http_class, *alert_class, *worker_class, *camera_class, *pushNotif_class, *cookie_class, *session_class, *local_class, *websocket_class, *navigation_class, *vibration_class, *notification_class, *location_class, *indexeddb_class;

  mrb = mrb_open();
  c = mrbc_context_new(mrb);

  if (!mrb)
  { /* handle error */
  }
  dom_class = mrb_define_class(mrb, "InternalDOM", mrb->object_class);
  mrb_define_class_method(
      mrb,
      dom_class,
      "add_event_listener",
      add_event_listener,
      MRB_ARGS_REQ(3));
  mrb_define_class_method(
      mrb,
      dom_class,
      "get_value",
      get_value,
      MRB_ARGS_REQ(2));
  // mrb_define_class_method(
  //     mrb,
  //     dom_class,
  //     "get_element",
  //     get_element,
  //     MRB_ARGS_REQ(2));

  http_class = mrb_define_class(mrb, "InternalHTTP", mrb->object_class);
  mrb_define_class_method(
      mrb,
      http_class,
      "http_request",
      http_request,
      MRB_ARGS_REQ(1));
  mrb_define_class_method(
      mrb,
      http_class,
      "http_post",
      http_post,
      MRB_ARGS_REQ(2));

  cookie_class = mrb_define_class(mrb, "InternalCookie", mrb->object_class);
  mrb_define_class_method(
      mrb,
      cookie_class,
      "set_cookie",
      set_cookie,
      MRB_ARGS_REQ(1));
  mrb_define_class_method(
      mrb,
      cookie_class,
      "get_cookie",
      get_cookie,
      MRB_ARGS_REQ(1));

  session_class = mrb_define_class(mrb, "InternalSession", mrb->object_class);
  mrb_define_class_method(
      mrb,
      session_class,
      "set_session",
      set_session,
      MRB_ARGS_REQ(2));
  mrb_define_class_method(
      mrb,
      session_class,
      "get_session",
      get_session,
      MRB_ARGS_REQ(2));

  local_class = mrb_define_class(mrb, "InternalLocal", mrb->object_class);
  mrb_define_class_method(
      mrb,
      local_class,
      "set_local",
      set_local,
      MRB_ARGS_REQ(2));
  mrb_define_class_method(
      mrb,
      local_class,
      "get_local",
      get_local,
      MRB_ARGS_REQ(2));

  websocket_class = mrb_define_class(mrb, "InternalWebSocket", mrb->object_class);
  mrb_define_class_method(
      mrb,
      websocket_class,
      "open_connection",
      open_connection,
      MRB_ARGS_REQ(2));

  navigation_class = mrb_define_class(mrb, "InternalNavigation", mrb->object_class);
  mrb_define_class_method(
      mrb,
      navigation_class,
      "navigate_to",
      navigate_to,
      MRB_ARGS_REQ(2));
  mrb_define_class_method(
      mrb,
      navigation_class,
      "get_param",
      get_param,
      MRB_ARGS_REQ(2));

  vibration_class = mrb_define_class(mrb, "InternalVibrationAPI", mrb->object_class);
  mrb_define_class_method(
      mrb,
      vibration_class,
      "vibrate_for",
      vibrate_for,
      MRB_ARGS_REQ(1));
  location_class = mrb_define_class(mrb, "InternalLocation", mrb->object_class);
  mrb_define_class_method(
      mrb,
      location_class,
      "get_location",
      get_location,
      MRB_ARGS_REQ(1));
  notification_class = mrb_define_class(mrb, "InternalNotificationAPI", mrb->object_class);
  mrb_define_class_method(
      mrb,
      notification_class,
      "send_notif",
      send_notif,
      MRB_ARGS_REQ(2));

  camera_class = mrb_define_class(mrb, "InternalCameraAPI", mrb->object_class);
  mrb_define_class_method(
      mrb,
      camera_class,
      "enable_camera",
      enable_camera,
      MRB_ARGS_REQ(1));

  pushNotif_class = mrb_define_class(mrb, "InternalPushNotifAPI", mrb->object_class);
  mrb_define_class_method(
      mrb,
      pushNotif_class,
      "send_pushNotif",
      send_pushNotif,
      MRB_ARGS_REQ(0));
  indexeddb_class = mrb_define_class(mrb, "InternalIndexedDB", mrb->object_class);
  mrb_define_class_method(
      mrb,
      indexeddb_class,
      "create_db",
      create_db,
      MRB_ARGS_REQ(2));
  mrb_define_class_method(
      mrb,
      indexeddb_class,
      "open_db",
      open_db,
      MRB_ARGS_REQ(3));
  alert_class = mrb_define_class(mrb, "InternalAlert", mrb->object_class);
  mrb_define_class_method(
      mrb,
      alert_class,
      "show_alert",
      show_alert,
      MRB_ARGS_REQ(1));
  worker_class = mrb_define_class(mrb, "InternalWorkerThread", mrb->object_class);
  mrb_define_class_method(
      mrb,
      worker_class,
      "start_worker",
      start_worker,
      MRB_ARGS_REQ(2));
  mrb_define_class_method(
      mrb,
      worker_class,
      "receive_msg",
      receive_msg,
      MRB_ARGS_REQ(2));
  return 0;
}

mrb_value load_file(char *name)
{
  mrb_value v;
  FILE *fp = fopen(name, "r");
  if (fp == NULL)
  {
    fprintf(stderr, "Cannot open file: %s\n", name);
    return mrb_nil_value();
  }
  printf("[Prism] Loading: %s\n", name);
  mrbc_filename(mrb, c, name);
  v = mrb_load_file_cxt(mrb, fp, c);
  fclose(fp);
  if (mrb->exc)
  {
    mrb_print_error(mrb);
    mrb->exc = NULL;
  }
  return v;
}

int load(char *main, char *config)
{
  const char *class_name;
  /*int i;
  for (i = 0; i < argc; i++) {
    FILE *lfp = fopen(argv[i], "rb");
    if (lfp == NULL) {
      printf("Cannot open library file: %s\n", argv[i]);
      mrbc_context_free(mrb, c);
      return;
    }
    mrb_load_file_cxt(mrb, lfp, c);
    fclose(lfp);
  }*/
  mrb_define_global_const(mrb, "JSON_CONFIG", mrb_str_new_cstr(mrb, config));
  load_file("src/prism.rb");
  app = load_file(main);
  struct RClass *prism_module = mrb_module_get(mrb, "Prism");
  struct RClass *mount_class = mrb_class_get_under(mrb, prism_module, "Mount");

  if (!mrb_obj_is_kind_of(mrb, app, mount_class))
  {
    class_name = mrb_obj_classname(mrb, app);

    fprintf(stderr, "[Prism] Error starting app.\n  Expected '%s' to return an instance of Prism::Mount but got a %s instead.\n  Did you remember to call Prism.mount on the last line?\n", main, class_name);

    return 1;
  }
  mrb_gc_register(mrb, app);

  return 0;
}

char *render()
{
  mrb_value result = mrb_funcall(mrb, app, "render", 0);
  if (mrb->exc)
  {
    mrb_print_error(mrb);
    mrb->exc = NULL;
  }
  return RSTRING_PTR(result);
}

void dispatch(char *message)
{
  mrb_value str = mrb_str_new_cstr(mrb, message);
  mrb_gc_register(mrb, str);
  mrb_funcall(mrb, app, "dispatch", 1, str);
  if (mrb->exc)
  {
    mrb_print_error(mrb);
    mrb->exc = NULL;
  }
  mrb_gc_unregister(mrb, str);
}

void event(char *message, char *id)
{
  mrb_value str = mrb_str_new_cstr(mrb, message);
  mrb_value str2 = mrb_str_new_cstr(mrb, id);
  mrb_funcall(mrb, app, "event", 2, str, str2);

  if (mrb->exc)
  {
    mrb_print_error(mrb);
    mrb->exc = NULL;
  }
}

void display_location(char *message, char *id)
{
  mrb_value str = mrb_str_new_cstr(mrb, message);
  mrb_value str2 = mrb_str_new_cstr(mrb, id);
  mrb_funcall(mrb, app, "display_location", 2, str, str2);

  if (mrb->exc)
  {
    mrb_print_error(mrb);
    mrb->exc = NULL;
  }
}

void display_value(char *text, char *id)
{
  mrb_value str = mrb_str_new_cstr(mrb, text);
  mrb_value str2 = mrb_str_new_cstr(mrb, id);

  mrb_funcall(mrb, app, "display_value", 2, str, str2);

  if (mrb->exc)
  {
    mrb_print_error(mrb);
    mrb->exc = NULL;
  }
}

void http_response(char *text, char *id)
{
  mrb_value str = mrb_str_new_cstr(mrb, text);
  mrb_value str2 = mrb_str_new_cstr(mrb, id);
  mrb_funcall(mrb, app, "http_response", 2, str, str2);

  if (mrb->exc)
  {
    mrb_print_error(mrb);
    mrb->exc = NULL;
  }
}

void session_value(char *text, char *id)
{
  mrb_value str = mrb_str_new_cstr(mrb, text);
  mrb_value str2 = mrb_str_new_cstr(mrb, id);

  mrb_funcall(mrb, app, "session_value", 2, str, str2);

  if (mrb->exc)
  {
    mrb_print_error(mrb);
    mrb->exc = NULL;
  }
}

void cookie_value(char *text, char *id)
{
  mrb_value str = mrb_str_new_cstr(mrb, text);
  mrb_value str2 = mrb_str_new_cstr(mrb, id);

  mrb_funcall(mrb, app, "cookie_value", 2, str, str2);

  if (mrb->exc)
  {
    mrb_print_error(mrb);
    mrb->exc = NULL;
  }
}

void local_value(char *text, char *id)
{
  mrb_value str = mrb_str_new_cstr(mrb, text);
  mrb_value str2 = mrb_str_new_cstr(mrb, id);

  mrb_funcall(mrb, app, "local_value", 2, str, str2);

  if (mrb->exc)
  {
    mrb_print_error(mrb);
    mrb->exc = NULL;
  }
}

void send_msg(char *text, char *id)
{
  mrb_value str = mrb_str_new_cstr(mrb, text);
  mrb_value str2 = mrb_str_new_cstr(mrb, id);

  mrb_funcall(mrb, app, "send_msg", 2, str, str2);

  if (mrb->exc)
  {
    mrb_print_error(mrb);
    mrb->exc = NULL;
  }
}

void param_value(char *text, char *id)
{
  mrb_value str = mrb_str_new_cstr(mrb, text);
  mrb_value str2 = mrb_str_new_cstr(mrb, id);

  mrb_funcall(mrb, app, "param_value", 2, str, str2);

  if (mrb->exc)
  {
    mrb_print_error(mrb);
    mrb->exc = NULL;
  }
}

void indexeddb_response(char *text, char *id)
{
  mrb_value str = mrb_str_new_cstr(mrb, text);
  mrb_value str2 = mrb_str_new_cstr(mrb, id);

  mrb_funcall(mrb, app, "indexeddb_response", 2, str, str2);

  if (mrb->exc)
  {
    mrb_print_error(mrb);
    mrb->exc = NULL;
  }
}

void worker_fn(char *text, char *id)
{
  mrb_value str = mrb_str_new_cstr(mrb, text);
  mrb_value str2 = mrb_str_new_cstr(mrb, id);

  mrb_funcall(mrb, app, "worker_fn", 2, str, str2);

  if (mrb->exc)
  {
    mrb_print_error(mrb);
    mrb->exc = NULL;
  }
}

void camera_output(char *text, char *id)
{
  mrb_value str = mrb_str_new_cstr(mrb, text);
  mrb_value str2 = mrb_str_new_cstr(mrb, id);

  mrb_funcall(mrb, app, "camera_output", 2, str, str2);

  if (mrb->exc)
  {
    mrb_print_error(mrb);
    mrb->exc = NULL;
  }
}
