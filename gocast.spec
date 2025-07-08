Name:           gocast
Version:        1.1.0
Release:        1%{?dist}
Summary:        A terminal-based weather forecast application

License:        MIT
URL:            https://github.com/example/gocast
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  golang >= 1.20
Requires:       ca-certificates

%description
Gocast is a simple terminal-based weather application that provides current weather
conditions and forecasts using the Open-Meteo API. It supports current conditions,
24-hour forecasts, and 7-day forecasts with ASCII art weather representations.
Now includes country specification support to disambiguate locations.

%prep
%setup -q

%build
# Binary is pre-built in the Makefile

%install
mkdir -p %{buildroot}/usr/local/bin
install -m 755 usr/local/bin/%{name} %{buildroot}/usr/local/bin/%{name}

%global debug_package %{nil}

%files
/usr/local/bin/%{name}

%changelog
* Mon Jun 30 2025 Package Maintainer <maintainer@example.com> - 1.1.0-1
- Added country specification support for location disambiguation
- Usage: gocast [options] <location> [country]
- Supports ISO country codes (GB, US, CA, etc.)
- Example: gocast fenton gb

* Mon Jun 30 2025 Package Maintainer <maintainer@example.com> - 1.0.0-1
- Initial RPM package for gocast weather application
- Includes current weather, 24h and 7d forecast functionality
- ASCII art weather representations