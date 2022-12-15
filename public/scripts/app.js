let geocoder;
let map;
let searchBox;

function init() {
    geocoder = new google.maps.Geocoder();

    map = new google.maps.Map(document.getElementById('map'));
    map.fitBounds({
        south: 45,
        west: -80,
        east: -63,
        north: 63
    });
    map.addListener('click', onClickMap);

    searchBox = new google.maps.places.Autocomplete(document.getElementById('search-input'), {
        bounds: {
            south: 45,
            west: -80,
            east: -63,
            north: 63
        },
        fields: ['address_components', 'geometry'],
        componentRestrictions: {
            country: 'ca'
        },
        strictBounds: true,
        types: ['route']
    });
    searchBox.addListener('place_changed', onSearch);
}

async function onSearch() {
    const info = searchBox.getPlace();
    display(undefined);
    if (!info || !info.formatted_address) return;

    const response = await axios.get('/api', {
        params: {
            route: parseName(info.address_components)
        }
    });

    if (response.status !== 200) {
        console.error('Failed to fetch API');
        return;
    }

    const data = {
        ...(response.data.route ?? {
            empty: true
        }),
        route: parseName(info.address_components)
    };

    display(data);

    if (!info.geometry || !info.geometry.location) return;
    const bounds = new google.maps.LatLngBounds();
    if (info.geometry.viewport) bounds.union(info.geometry.viewport);
    else bounds.extend(info.geometry.location);
    map.fitBounds(bounds);
}

function onClickMap(event) {
    const request = {
        location: event.latLng
    };

    display(undefined);
    geocoder.geocode(request, async (results, status) => {
        if (status !== 'OK') {
            console.error(`Failed to geocode: ${status}`);
            return;
        }

        const result = results.find(elem => elem.types.includes('route'));
        if (!result) {
            console.info('User clicked on non-route, ignoring');
            return;
        }

        const id = result.place_id;
        const response = await axios.get('/api', {
            params: {
                placeId: id
            }
        });

        if (response.status !== 200) {
            console.error('Failed to fetch API');
            return;
        }

        const data = {
            ...(response.data.route ?? {
                empty: true
            }),
            route: parseName(result.address_components)
        };

        display(data);
    });
}

function parseName(components) {
    const route = components.find(elem => elem.types.includes('route'));
    const locality = components.find(elem => elem.types.includes('locality'));

    return `${route.long_name}, ${locality.long_name}`;
}

function display(data) {
    const container = $('#route');
    container.empty();
    const title = `<h2>${data ? data.route : 'Aucune sélection'}</h2>`;
    container.append(title);
    if (data) {
        if (data.rank !== undefined) {
            container.append(`<div class="row">
                <h6 class="col-3 header">Rang global</h6>
                <h6 class="col-9 text-end">#${data.rank}/${data.maxRank}</h6>
            </div>`);
        }

        if (data.global_value !== undefined) {
            container.append(`<div class="row">
                <h6 class="col-3 header">Dommages globaux</h6>
                <div class="progress col-9">
                    <div class="progress-bar damage-bar" style="width: ${data.global_value}%"></div>
                </div>
            </div>`);
        }

        if (data.local_value !== undefined) {
            container.append(`<div class="row">
                <h6 class="col-3 header">Dommages locaux</h6>
                <div class="progress col-9">
                    <div class="progress-bar damage-bar" style="width: ${data.local_value}%"></div>
                </div>
            </div>`);
        }

        if (data.empty) {
            container.append('<h5>Aucune données disponibles...</h5>');
        }
    }
}