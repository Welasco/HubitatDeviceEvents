import type { Component } from 'solid-js';
import { onMount } from 'solid-js';
import { createSignal, createResource, For } from 'solid-js';
import TimePicker from "@rnwonder/solid-date-picker/timePicker";
import DatePicker, {
  PickerValue,
  TimeValue,
} from '@rnwonder/solid-date-picker';

import { CalendarIcon } from './Icons/Calendar';
import { NextIcon } from './Icons/Next';
import { PrevIcon } from './Icons/Prev';

const apiURL = 'http://vwsstorage.localdomain:3000/api/v1';
const apiGetDevices = apiURL + '/device';
const apiGetDevice = apiURL + '/device/{id}';
const apiGetDeviceEvents = apiURL + '/device/{id}/event?start={start}&end={end}';

interface IDevice {
  Id: string,
  name: string,
  label: string,
  type: string,
  room: string
}

const Device = {
  Id: '',
  name: '',
  label: '',
  type: '',
  room: ''
}

const fetchDevices = async () =>
  (await fetch(`http://vwsstorage.localdomain:3000/api/v1/device`)).json();

async function fetchAPI(apiUri: string) {
  console.log('fetchAPI: ', apiUri)
  //return (await fetch(apiUri)).json();
  try {
    const response = (await fetch(apiUri)).json();
    return response
  } catch (error) {
    console.error('Error fetching data:', error);
    throw error
  }
}

// function fetchAPI(apiUri: string) {
//   console.log('fetchAPI: ', apiUri)
//   try {
//     const response = async () => (await fetch(apiUri)).json();
//     return response
//   } catch (error) {
//     console.error('Error fetching data:', error);
//     throw error
//   }
// }
// const fetchDevices = async () =>
//   (await fetch(`http://vwsstorage.localdomain:3000/api/v1/device`)).json();

const App: Component = () => {
  //const [area, setArea] = createSignal('device_details')
  const [area, setArea] = createSignal('device_events')
  const [readDevice, setReadDevice] = createSignal<number>(1);
  const toggleDevice = () => setReadDevice(readDevice() == 1 ? setReadDevice(2) : setReadDevice(1))
  const [devices] = createResource(readDevice, async () => await fetchAPI(apiGetDevices));

  const [readDeviceDetail, setReadDeviceDetail] = createSignal<IDevice>(Device);
  const [deviceDetail] = createResource(readDeviceDetail, async () => await fetchAPI(apiGetDevice.replace('{id}', readDeviceDetail().Id)));

  const [toggleReadDeviceEvents, settoggleReadDeviceEvents] = createSignal<number>(1);
  const toggleDeviceEvents = () => toggleReadDeviceEvents() == 1 ? settoggleReadDeviceEvents(2) : settoggleReadDeviceEvents(1)
  const [deviceEventsToggle] = createResource(toggleReadDeviceEvents, async () => await fetchAPI(apiGetDeviceEvents.replace('{id}', readDeviceDetail().Id).replace('{start}', buildDateTime('from')).replace('{end}', buildDateTime('to'))));
  //const [readDeviceEvents, setReadDeviceEvents] = createSignal<IDevice>(Device);
  //const [deviceEvents] = createResource(readDeviceEvents, async () => await fetchAPI(apiGetDeviceEvents.replace('{id}', readDeviceDetail().Id).replace('{start}', buildDateTime('from')).replace('{end}', buildDateTime('to'))));

  const [singleDateFrom, setSingleDateFrom] = createSignal<PickerValue>({
    label: (new Date(buildDateTime("initFromDate")).toLocaleString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })),
    value: {
      selected: buildDateTime("initFromDate")
    },
  });
  const [ldeDF] = createResource(singleDateFrom, loadDeviceEvents)

  const [timeValueFrom, setTimeValueFrom] = createSignal<TimeValue>({
    value: {
      hour: parseInt(buildDateTime("initFromHour")),
      minute: new Date().getMinutes(),
      second: new Date().getSeconds()
    },
    label: (new Date(buildDateTime("initFromDate")).toLocaleTimeString()),
  });
  const [ldeTF] = createResource(timeValueFrom, loadDeviceEvents)

  const [singleDateTo, setSingleDateTo] = createSignal<PickerValue>({
    label: (new Date().toLocaleString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })),
    value: {
      selected: buildDateTime("current")
    },
  });
  const [ldeDT] = createResource(singleDateTo, loadDeviceEvents)

  const [timeValueTo, setTimeValueTo] = createSignal<TimeValue>({
    value: {
      hour: new Date().getHours(),
      minute: new Date().getMinutes(),
      second: new Date().getSeconds()
    },
    label: (new Date().toLocaleTimeString()),
  });
  const [ldeTT] = createResource(timeValueTo, loadDeviceEvents)


  function loadDevice(device: IDevice) {
    console.log('loadDevice: ', device)
    setReadDeviceDetail(device)
    //setReadDeviceEvents(device)
    loadDeviceEvents()
  }

  function loadDeviceEvents() {
    console.log('loadDeviceEvents: ', readDeviceDetail)
    //setReadDeviceEvents(readDeviceDetail())
    toggleDeviceEvents()
  }

  // function loadDeviceDetails(device: IDevice) {
  //   console.log('loadDeviceDetails: ', device)
  //   //setReadDeviceDetail({ ...readDeviceDetail(), Id: device.Id, name: device.name, label: device.label, type: device.type, room: Math.random().toString()})
  //   setReadDeviceDetail(device)
  //   console.log('readDeviceDetail: ', readDeviceDetail())
  //   console.log('test replace: ', apiGetDevice.replace('{id}', readDeviceDetail().Id))

  // }

  // function loadDeviceEvents(device: IDevice) {
  //   console.log('loadDeviceEvents: ', device)
  //   setReadDeviceDetail(device)
  // }

  function buildDateTime(fromOrtoOrCurrent: string): string {

    function fixZeroDateTime(dateTime: string): string {
      let receivedDateTime = new Date(dateTime)
      let hours = receivedDateTime.getHours()
      let minutes = receivedDateTime.getMinutes()
      let seconds = receivedDateTime.getSeconds()

      let currenthours = ("0" + hours).slice(-2);
      let currentminutes = ("0" + minutes).slice(-2);
      let currentseconds = ("0" + seconds).slice(-2);

      let returnFixedZeroDateTime = receivedDateTime.toLocaleDateString('sv-SE')+"T"+currenthours+":"+currentminutes+":"+currentseconds+".000"
      //let returnFixedZeroDateTime = receivedDateTime.toLocaleDateString('sv-SE')+"T"+currenthours+":"+currentminutes+":"+currentseconds+".000Z"

      return returnFixedZeroDateTime
    }

    if (fromOrtoOrCurrent == 'from') {
      let fromDateTime = singleDateFrom().value?.selected?.split('T')[0] + 'T' + timeValueFrom().value.hour + ':' + timeValueFrom().value.minute + ':' + timeValueFrom().value.second
      console.log('fromDateTime: ', fromDateTime)
      //let fromDateTime = singleDateFrom().value.selected.split('T')[0] + 'T' + timeValueFrom().value.hour + ':' + timeValueFrom().value.minute + ':' + timeValueFrom().value.second
      return fromDateTime
    } else if (fromOrtoOrCurrent == 'to'){
      let toDateTime = singleDateTo().value?.selected?.split('T')[0] + 'T' + timeValueTo().value.hour + ':' + timeValueTo().value.minute + ':' + timeValueTo().value.second
      console.log('toDateTime: ', toDateTime)
      return toDateTime
    } else if(fromOrtoOrCurrent == 'current'){
      let currentDateTime = new Date()
      let returnCurrentDateTime = fixZeroDateTime(currentDateTime.toLocaleString())
      console.log('returnCurrentDateTime: ', returnCurrentDateTime)
      return returnCurrentDateTime
    } else if(fromOrtoOrCurrent == 'initFromDate'){
      let initFromDateTime = new Date()
      initFromDateTime.setHours(initFromDateTime.getHours() - 8)

      let initFromDateTimeStr = fixZeroDateTime(initFromDateTime.toLocaleString())
      console.log('initFromDateTime: ', initFromDateTimeStr)
      return initFromDateTimeStr
    } else {
      let initFromHour = new Date()
      initFromHour.setHours(initFromHour.getHours() - 8)
      let returnInitFromHour = initFromHour.getHours()
      console.log('initFromHour: ', returnInitFromHour)
      return returnInitFromHour.toString()
    }
  }

  toggleDevice()

  return (
  <>
    <header class="bg-slate-950 text-white p-3 text-2xl font-bold h-[60px]">HubitatDeviceEvents</header>
    <div class="flex flex-row h-[calc(100vh-100px)] bg-white">
      <aside class="bg-slate-200 p-2 w-1/4 overflow-auto">
        <section class="h-[40px] bg-gray-300 px-2 space-x-2 flex flex-row"></section>
        <div class="bg-slate-600 h-[calc(100vh-160px)] px-2">
          <p>Select a device:</p>

          <select id="deviceList" size="40" class="px-2 w-full">
            <For each={devices()}>
              {(device) => (
                <option value={device.id} onClick={() => loadDevice(device)}>{device.label}</option>
              )}
            </For>
          </select>

          <br></br>
          <button onClick={toggleDevice} class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 border border-blue-700 rounded">Load Devices</button>
          <button onClick={toggleDeviceEvents} class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 border border-blue-700 rounded">Load Device Events</button>

          <br></br>
          {/* <span>{devices.loading && "Loading..."}</span>
            <pre>{JSON.stringify(devices(), null, 2)}</pre> */}
        </div>
      </aside >
      <main id="main-area" class="h-full w-full overflow-auto flex flex-row flex-wrap p-2">
        <nav class="flex border-b">
          <button class={(area() == 'device_details' ? 'bg-white inline-block border-l border-t border-r rounded-t py-2 px-4 text-blue-700 font-semibold' : 'bg-white inline-block py-2 px-4 text-blue-500 hover:text-blue-800 font-semibold')} onclick={() => setArea('device_details')}>
            Device Details
          </button>
          <button class={(area() == 'device_events' ? 'bg-white inline-block border-l border-t border-r rounded-t py-2 px-4 text-blue-700 font-semibold' : 'bg-white inline-block py-2 px-4 text-blue-500 hover:text-blue-800 font-semibold')} onclick={() => setArea('device_events')}>
            Device Events
          </button>
        </nav>

        <div id='device_details' class={'w-full md:w-full bg-blue-200 h-[calc(100vh-160px)] text-black ' + (area() == 'device_details' ? '' : 'hidden')}>
          <p class='text-lg'>Device Details</p>
          <p class='text-lg'>Test readDeviceDetail() trigger: {readDeviceDetail().Id}</p>
          test333333333333333333333333333333
          test44444
          <p class='text-lg'>{readDeviceDetail() == Device ? '' : JSON.stringify(deviceDetail(), null, 2)}</p>
        </div>
        <div id='device_events' class={'w-full md:w-full  bg-blue-200 h-[calc(100vh-160px)] text-black ' + (area() == 'device_events' ? '' : 'hidden')}>

          {/* ######################################################################################## */}
          {/* ############################### Calendar Component ##################################### */}
          <div id="date-time-menu" class='bg-red-500 h-[60px] flex flex-row flex-wrap p-3 space-x-2'>
            <div id="date-time-from" class="text-white flex flex-row flex-wrap p-3 space-x-2">
              <p class={'text-lg '}>From:</p>
              <div>
                <DatePicker
                  removeNavButtons
                  weekDaysType="single"
                  shouldHighlightWeekends
                  hideOutSideDays
                  afterNextButtonAreaJSX={({ handleNextMonth, handlePrevMonth }) => (
                    <div class="">
                      <button onClick={handlePrevMonth}>
                        <PrevIcon />
                      </button>
                      <button onClick={handleNextMonth}>
                        <NextIcon />
                      </button>
                    </div>
                  )}
                  renderInput={({ showDate, value }) => (
                    <div class="custom-input">
                      {/* <input value={value().label} readOnly /> */}
                      <input value={singleDateFrom().label} readOnly />
                      <button onClick={showDate}>
                        <CalendarIcon />
                      </button>
                    </div>
                  )}
                  onChange={() => {toggleDeviceEvents}}
                  value={singleDateFrom}
                  setValue={setSingleDateFrom}
                  type={'single'}
                  monthSelectorFormat={'long'}
                />
              </div>
              <div>
                <TimePicker
                  value={timeValueFrom}
                  setValue={setTimeValueFrom}
                  shouldCloseOnSelect
                  allowedView={['hour', 'minute', 'second']}
                  renderInput={({ showTime, value }) => (
                    <div class="custom-input">
                      <input value={timeValueFrom().label} readOnly />
                      <button onClick={showTime}>
                        <CalendarIcon />
                      </button>
                    </div>
                  )}
                  timeAnalogClockCenterDotClass='text-red-500'
                  timeAnalogNumberClass='text-black'
                  timeAnalogClockHandClass='text-black'
                  timeAnalogWrapperClass='text-black'
                  timePickerBottomAreaClass='text-black'
                  timePickerMeridiemBtnClass='text-black'
                  timePickerTopAreaClass='text-black'
                  timePickerWrapperClass='text-black'
                />
              </div>
            </div>
            <div id="date-time-to" class="text-white flex flex-row flex-wrap p-3 space-x-2">
              <p class={'text-lg '}>To:</p>
              <div>
                <DatePicker
                  removeNavButtons
                  weekDaysType="single"
                  shouldHighlightWeekends
                  hideOutSideDays
                  afterNextButtonAreaJSX={({ handleNextMonth, handlePrevMonth }) => (
                    <div class="">
                      <button onClick={handlePrevMonth}>
                        <PrevIcon />
                      </button>
                      <button onClick={handleNextMonth}>
                        <NextIcon />
                      </button>
                    </div>
                  )}
                  renderInput={({ showDate, value }) => (
                    <div class="custom-input">
                      <input value={singleDateTo().label} readOnly />
                      <button onClick={showDate}>
                        <CalendarIcon />
                      </button>
                    </div>
                  )}
                  onChange={() => {toggleDeviceEvents}}
                  value={singleDateTo}
                  setValue={setSingleDateTo}
                  type={'single'}
                  monthSelectorFormat={'long'}
                />
              </div>

              <div>
                <TimePicker
                  value={timeValueTo}
                  setValue={setTimeValueTo}
                  shouldCloseOnSelect
                  allowedView={['hour', 'minute', 'second']}
                  renderInput={({ showTime, value }) => (
                    <div class="custom-input">
                      <input value={timeValueTo().label} readOnly />
                      <button onClick={showTime}>
                        <CalendarIcon />
                      </button>
                    </div>
                  )}
                  onChange={() => {toggleDeviceEvents}}
                  timeAnalogClockCenterDotClass='text-red-500'
                  timeAnalogNumberClass='text-black'
                  timeAnalogClockHandClass='text-black'
                  timeAnalogWrapperClass='text-black'
                  timePickerBottomAreaClass='text-black'
                  timePickerMeridiemBtnClass='text-black'
                  timePickerTopAreaClass='text-black'
                  timePickerWrapperClass='text-black'
                />
              </div>
            </div>
          </div>
          {/* ######################################################################################## */}
          <p class={'text-lg '}>Device Events:</p>
          deviceEvents
          {/* <p class='text-lg'>{readDeviceDetail() == Device ? '' : JSON.stringify(deviceEvents(), null, 2)}</p> */}
          <p class='text-lg'>{JSON.stringify(deviceEventsToggle(), null, 2)}</p>
          device_events
          <p class='text-lg'>Value of toggleReadDeviceEvents: {toggleReadDeviceEvents()}</p>
          <p class='text-lg'>Value of string toggleReadDeviceEvents: {toggleReadDeviceEvents().toString()}</p>
        </div>
      </main>
    </div >
    <footer class="h-[40px] text-white bg-black text-2xl px-2">footer</footer>
  </>
);
};

export default App;
