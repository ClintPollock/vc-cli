/*
 * VERACODE.SOURCE.CONFIDENTIAL.security-logging-java.b5f7196e3c8d51cae90d11c2f37240654e19bcc09da964086d43ce67f1f200de
 *
 * Copyright Veracode Inc., 2017
 */
package com.veracode.security.logging;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.slf4j.MDC;
import org.slf4j.Marker;

import java.util.Map;

/**
 * SecureLogger
 */
public class SecureLogger implements Logger {

    private final Logger baseLogger;

    public static SecureLogger getLogger(String name) {
        return new SecureLogger(name);
    }

    public static SecureLogger getLogger(Class clazz) {
        return new SecureLogger(clazz);
    }

    protected SecureLogger(String name) {
        baseLogger = LoggerFactory.getLogger(name);
    }

    protected SecureLogger(Class clazz) {
        baseLogger = LoggerFactory.getLogger(clazz);
    }

    public Logger getBaseLogger() {
        return baseLogger;
    }

    @Override
    public String getName() {
        return baseLogger.getName();
    }

    @Override
    public boolean isTraceEnabled() {
        return baseLogger.isTraceEnabled();
    }

    @Override
    public boolean isTraceEnabled(Marker marker) {
        return baseLogger.isTraceEnabled(marker);
    }

    @Override
    public void trace(String msg) {
        String escapedString = SecureLoggerUtil.escapeMessage(msg);
        baseLogger.trace(escapedString);
    }

    @Override
    public void trace(String format, Object arg) {
        Object escapedString = SecureLoggerUtil.escapeMessage(arg);
        baseLogger.trace(format, escapedString);
    }

    @Override
    public void trace(String format, Object arg1, Object arg2) {
        Object str1 = SecureLoggerUtil.escapeMessage(arg1);
        Object str2 = SecureLoggerUtil.escapeMessage(arg2);
        baseLogger.trace(format, str1, str2);
    }

    @Override
    public void trace(String format, Object... arguments) {
        baseLogger.trace(format, SecureLoggerUtil.escapeMessages(arguments));
    }

    @Override
    public void trace(String msg, Throwable t) {
        String escapedString = SecureLoggerUtil.escapeMessage(msg);
        baseLogger.trace(escapedString, new SecureExceptionWrapper(t));
    }

    @Override
    public void trace(Marker marker, String msg) {
        String escapedString = SecureLoggerUtil.escapeMessage(msg);
        baseLogger.trace(marker, escapedString);
    }

    @Override
    public void trace(Marker marker, String format, Object arg) {
        baseLogger.trace(marker, format, SecureLoggerUtil.escapeMessages(arg));
    }

    @Override
    public void trace(Marker marker, String format, Object arg1, Object arg2) {
        Object str1 = SecureLoggerUtil.escapeMessage(arg1);
        Object str2 = SecureLoggerUtil.escapeMessage(arg2);
        baseLogger.trace(marker, format, str1, str2);
    }

    @Override
    public void trace(Marker marker, String format, Object... argArray) {
        baseLogger.trace(marker, format, SecureLoggerUtil.escapeMessages(argArray));
    }

    @Override
    public void trace(Marker marker, String msg, Throwable t) {
        String escapedString = SecureLoggerUtil.escapeMessage(msg);
        baseLogger.trace(marker, escapedString, new SecureExceptionWrapper(t));
    }

    @Override
    public boolean isDebugEnabled() {
        return baseLogger.isDebugEnabled();
    }

    @Override
    public boolean isDebugEnabled(Marker marker) {
        return baseLogger.isDebugEnabled(marker);
    }

    @Override
    public void debug(String msg) {
        String escapedString = SecureLoggerUtil.escapeMessage(msg);
        baseLogger.debug(escapedString);
    }

    @Override
    public void debug(String format, Object arg) {
        baseLogger.debug(format, SecureLoggerUtil.escapeMessage(arg));
    }

    @Override
    public void debug(String format, Object arg1, Object arg2) {
        Object str1 = SecureLoggerUtil.escapeMessage(arg1);
        Object str2 = SecureLoggerUtil.escapeMessage(arg2);
        baseLogger.debug(format, str1, str2);
    }

    @Override
    public void debug(String format, Object... arguments) {
        baseLogger.debug(format, SecureLoggerUtil.escapeMessages(arguments));
    }

    @Override
    public void debug(String msg, Throwable t) {
        String escapedString = SecureLoggerUtil.escapeMessage(msg);
        baseLogger.debug(escapedString, new SecureExceptionWrapper(t));
    }

    @Override
    public void debug(Marker marker, String msg) {
        String escapedString = SecureLoggerUtil.escapeMessage(msg);
        baseLogger.debug(marker, escapedString);
    }

    @Override
    public void debug(Marker marker, String format, Object arg) {
        baseLogger.debug(marker, format, SecureLoggerUtil.escapeMessage(arg));
    }

    @Override
    public void debug(Marker marker, String format, Object arg1, Object arg2) {
        Object str1 = SecureLoggerUtil.escapeMessage(arg1);
        Object str2 = SecureLoggerUtil.escapeMessage(arg2);
        baseLogger.debug(marker, format, str1, str2);
    }

    @Override
    public void debug(Marker marker, String format, Object... arguments) {
        baseLogger.debug(marker, format, SecureLoggerUtil.escapeMessages(arguments));
    }

    @Override
    public void debug(Marker marker, String msg, Throwable t) {
        String escapedString = SecureLoggerUtil.escapeMessage(msg);
        baseLogger.debug(marker, escapedString, new SecureExceptionWrapper(t));
    }

    @Override
    public boolean isInfoEnabled() {
        return baseLogger.isInfoEnabled();
    }

    @Override
    public boolean isInfoEnabled(Marker marker) {
        return baseLogger.isInfoEnabled(marker);
    }

    @Override
    public void info(String msg) {
        String escapedString = SecureLoggerUtil.escapeMessage(msg);
        baseLogger.info(escapedString);
    }

    @Override
    public void info(String format, Object arg) {
        baseLogger.info(format, SecureLoggerUtil.escapeMessage(arg));
    }

    @Override
    public void info(String format, Object arg1, Object arg2) {
        Object str1 = SecureLoggerUtil.escapeMessage(arg1);
        Object str2 = SecureLoggerUtil.escapeMessage(arg2);
        baseLogger.info(format, str1, str2);
    }

    @Override
    public void info(String format, Object... arguments) {
        baseLogger.info(format, SecureLoggerUtil.escapeMessages(arguments));
    }

    @Override
    public void info(String msg, Throwable t) {
        String escapedString = SecureLoggerUtil.escapeMessage(msg);
        baseLogger.info(escapedString, new SecureExceptionWrapper(t));
    }

    @Override
    public void info(Marker marker, String msg) {
        String escapedString = SecureLoggerUtil.escapeMessage(msg);
        baseLogger.info(marker, escapedString);
    }

    @Override
    public void info(Marker marker, String format, Object arg) {
        baseLogger.info(marker, format, SecureLoggerUtil.escapeMessage(arg));
    }

    @Override
    public void info(Marker marker, String format, Object arg1, Object arg2) {
        Object str1 = SecureLoggerUtil.escapeMessage(arg1);
        Object str2 = SecureLoggerUtil.escapeMessage(arg2);
        baseLogger.info(marker, format, str1, str2);
    }

    @Override
    public void info(Marker marker, String format, Object... arguments) {
        baseLogger.info(marker, format, SecureLoggerUtil.escapeMessages(arguments));
    }

    @Override
    public void info(Marker marker, String msg, Throwable t) {
        String escapedString = SecureLoggerUtil.escapeMessage(msg);
        baseLogger.info(marker, escapedString, new SecureExceptionWrapper(t));
    }

    @Override
    public boolean isWarnEnabled() {
        return baseLogger.isWarnEnabled();
    }

    @Override
    public boolean isWarnEnabled(Marker marker) {
        return baseLogger.isWarnEnabled(marker);
    }

    @Override
    public void warn(String msg) {
        String escapedString = SecureLoggerUtil.escapeMessage(msg);
        baseLogger.warn(escapedString);
    }

    @Override
    public void warn(String format, Object arg) {
        baseLogger.warn(format, SecureLoggerUtil.escapeMessage(arg));
    }

    @Override
    public void warn(String format, Object... arguments) {
        baseLogger.warn(format, SecureLoggerUtil.escapeMessages(arguments));
    }

    @Override
    public void warn(String format, Object arg1, Object arg2) {
        Object str1 = SecureLoggerUtil.escapeMessage(arg1);
        Object str2 = SecureLoggerUtil.escapeMessage(arg2);
        baseLogger.warn(format, str1, str2);
    }

    @Override
    public void warn(String msg, Throwable t) {
        String escapedString = SecureLoggerUtil.escapeMessage(msg);
        baseLogger.warn(escapedString, new SecureExceptionWrapper(t));
    }

    @Override
    public void warn(Marker marker, String msg) {
        String escapedString = SecureLoggerUtil.escapeMessage(msg);
        baseLogger.warn(marker, escapedString);
    }

    @Override
    public void warn(Marker marker, String format, Object arg) {
        baseLogger.warn(marker, format, SecureLoggerUtil.escapeMessage(arg));
    }

    @Override
    public void warn(Marker marker, String format, Object arg1, Object arg2) {
        Object str1 = SecureLoggerUtil.escapeMessage(arg1);
        Object str2 = SecureLoggerUtil.escapeMessage(arg2);
        baseLogger.warn(marker, format, str1, str2);
    }

    @Override
    public void warn(Marker marker, String format, Object... arguments) {
        baseLogger.warn(marker, format, SecureLoggerUtil.escapeMessages(arguments));
    }

    @Override
    public void warn(Marker marker, String msg, Throwable t) {
        String escapedString = SecureLoggerUtil.escapeMessage(msg);
        baseLogger.warn(marker, escapedString, new SecureExceptionWrapper(t));
    }

    @Override
    public boolean isErrorEnabled() {
        return baseLogger.isErrorEnabled();
    }

    @Override
    public boolean isErrorEnabled(Marker marker) {
        return baseLogger.isErrorEnabled(marker);
    }

    @Override
    public void error(String msg) {
        String escapedString = SecureLoggerUtil.escapeMessage(msg);
        baseLogger.error(escapedString);
    }

    @Override
    public void error(String format, Object arg) {
        baseLogger.error(format, SecureLoggerUtil.escapeMessage(arg));
    }

    @Override
    public void error(String format, Object arg1, Object arg2) {
        Object str1 = SecureLoggerUtil.escapeMessage(arg1);
        Object str2 = SecureLoggerUtil.escapeMessage(arg2);
        baseLogger.error(format, str1, str2);
    }

    @Override
    public void error(String format, Object... arguments) {
        baseLogger.error(format, SecureLoggerUtil.escapeMessages(arguments));
    }

    @Override
    public void error(String msg, Throwable t) {
        String escapedString = SecureLoggerUtil.escapeMessage(msg);
        baseLogger.error(escapedString, new SecureExceptionWrapper(t));
    }

    @Override
    public void error(Marker marker, String msg) {
        String escapedString = SecureLoggerUtil.escapeMessage(msg);
        baseLogger.error(marker, escapedString);
    }

    @Override
    public void error(Marker marker, String format, Object arg) {
        baseLogger.error(marker, format, SecureLoggerUtil.escapeMessage(arg));
    }

    @Override
    public void error(Marker marker, String format, Object arg1, Object arg2) {
        Object str1 = SecureLoggerUtil.escapeMessage(arg1);
        Object str2 = SecureLoggerUtil.escapeMessage(arg2);
        baseLogger.error(marker, format, str1, str2);
    }

    @Override
    public void error(Marker marker, String format, Object... arguments) {
        baseLogger.error(marker, format, SecureLoggerUtil.escapeMessages(arguments));
    }

    @Override
    public void error(Marker marker, String msg, Throwable t) {
        String escapedString = SecureLoggerUtil.escapeMessage(msg);
        baseLogger.error(marker, escapedString, new SecureExceptionWrapper(t));
    }

    public static void addContext(Map<String, String> context) {
        if (context != null) {
            for (Map.Entry<String, String> entry : context.entrySet()) {
                addContext(entry.getKey(), entry.getValue());
            }
        }
    }

    public static void addContextIfEmpty(String key, String value) {
        if (key != null && MDC.get(SecureLoggerUtil.escapeMessage(key)) == null) {
            addContext(key, value);
        }
    }

    public static void addContext(String key, String value) {
        if (key != null) {
            MDC.put(SecureLoggerUtil.escapeMessage(key), SecureLoggerUtil.escapeMessage(value));
        }
    }

    public static void removeContext(Map<String, String> context) {
        if (context != null) {
            for (String key : context.keySet()) {
                removeContext(key);
            }
        }
    }

    public static void removeContext(String key) {
        if (key != null) {
            MDC.remove(SecureLoggerUtil.escapeMessage(key));
        }
    }

    public static void clearContext() {
        MDC.clear();
    }
}
